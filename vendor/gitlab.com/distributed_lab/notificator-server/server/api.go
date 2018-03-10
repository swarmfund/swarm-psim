package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type APIService struct {
	router            *web.Mux
	requestDispatcher *RequestDispatcher
	requests          conf.RequestsConf
	log               *logrus.Entry
}

func NewAPIService(log *logrus.Logger, requests conf.RequestsConf) *APIService {
	return &APIService{
		router:            web.New(),
		requestDispatcher: NewRequestDispatcher(),
		requests:          requests,
		log:               log.WithField("service", "api"),
	}
}

func (service *APIService) Init(cfg conf.Config) {
	r := service.router

	// middleware
	r.Use(ContentTypeMiddleware("application/json"))
	r.Use(LogMiddleware(service.log))
	r.Use(CheckAuthMiddleware(cfg.HTTP().AllowUntrusted))
	// routes

	r.Post("/", service.rootHandler)
}

func (service *APIService) Run(cfg conf.Config) {
	service.router.Compile()

	http.Handle("/", service.router)

	httpConf := cfg.HTTP()
	addr := fmt.Sprintf("%s:%d", httpConf.Host, httpConf.Port)

	service.log.WithField("starting on ", addr).Info()
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		service.log.WithError(err).Fatal("terminated")
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

	apiRequest := new(types.APIRequest)
	err = json.Unmarshal(body, apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"reason": "invalid request body", "msg": "%s"}`, err)
		return
	}
	result, err := service.requestDispatcher.Dispatch(apiRequest, service.requests)
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
