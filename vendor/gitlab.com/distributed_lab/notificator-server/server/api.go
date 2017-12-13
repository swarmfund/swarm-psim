package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"gitlab.com/distributed_lab/notificator-server/auth"
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/log"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type APIService struct {
	router            *web.Mux
	requestDispatcher *RequestDispatcher
	log               *logrus.Entry
}

func NewAPIService() *APIService {
	return &APIService{
		router:            web.New(),
		requestDispatcher: NewRequestDispatcher(),
		log:               log.WithField("service", "api"),
	}
}

func (service *APIService) Init() {
	r := service.router

	// middleware

	r.Use(ContentTypeMiddleware("application/json"))
	r.Use(LogMiddleware(service.log))

	// routes

	r.Post("/", service.rootHandler)

}

func (service *APIService) Run() {
	httpConf := conf.GetHTTPConf()
	service.router.Compile()

	http.Handle("/", service.router)

	addr := fmt.Sprintf("%s:%d", httpConf.Host, httpConf.Port)

	service.log.Infof("starting on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		service.log.WithError(err).Fatal("terminated")
	}
}

func checkAuth(r *http.Request, body []byte) (bool, error) {
	authorization := r.Header.Get("authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		// check just token
		key := strings.TrimPrefix(authorization, "Bearer ")
		pair, err := q.NewQ().Auth().ByPublic(key)
		if err != nil {
			return false, err
		}
		return pair != nil, nil
	} else {
		// legacy signature
		signature := r.Header.Get("x-signature")
		pair, err := q.NewQ().Auth().ByPublic(authorization)
		if err != nil {
			return false, err
		}
		return !(pair == nil || !auth.Verify(pair, body, signature)), nil
	}
}

func (service *APIService) rootHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"reason": "internal server error", "msg": "%s"}`, err)
		return
	}

	if !conf.GetHTTPConf().AllowUntrusted {
		ok, err := checkAuth(r, body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"reason": "internal server error", "msg": "%s"}`, err)
			return
		}

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"reason": "signature mismatch"}`)
			return
		}
	}

	apiRequest := new(types.APIRequest)
	err = json.Unmarshal(body, apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"reason": "invalid request body", "msg": "%s"}`, err)
		return
	}
	result, err := service.requestDispatcher.Dispatch(apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"reason": "internal server error", "msg": "%s"}`, err)
		return
	}

	switch result.Type {
	case DispatchResultSuccess:
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{}`)
	case DispatchResultLimitExceeded:
		w.WriteHeader(http.StatusTooManyRequests)
		response := APIErrorsResponse{
			Errors: []APIErrorResponse{
				{
					Status:      http.StatusTooManyRequests,
					IsPermanent: result.IsPermanent,
					RetryIn:     NewAPIDuration(result.RetryIn),
				},
			},
		}
		body, err := json.Marshal(&response)
		if err != nil {
			panic(err)
		}
		w.Write(body)
	case DispatchResultUnknownType:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{}`)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{}`)
	}
}

type APIErrorsResponse struct {
	Errors []APIErrorResponse `json:"errors"`
}

type APIErrorResponse struct {
	Status      int          `json:"status"`
	IsPermanent bool         `json:"is_permanent"`
	RetryIn     *APIDuration `json:"retry_in,omitempty"`
}

type APIDuration struct {
	*time.Duration
}

func NewAPIDuration(duration *time.Duration) *APIDuration {
	if duration == nil {
		return nil
	}
	return &APIDuration{Duration: duration}
}
func (d *APIDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(d.Seconds()))
}
