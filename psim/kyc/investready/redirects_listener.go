package investready

import (
	"context"
	"net/http"
	"time"

	"fmt"

	"encoding/json"

	"io/ioutil"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

type redirectedRequest struct {
	AccountID string `json:"account_id"`
	OauthCode string `json:"oauth_code"`
}

func (r redirectedRequest) Validate() (validationErr string) {
	if r.AccountID == "" {
		return "account_id cannot be empty."
	}
	if r.OauthCode == "" {
		return "oauth_code cannot be empty."
	}

	return ""
}

type RedirectsListener struct {
	log         *logan.Entry
	config      RedirectsConfig
	investReady InvestReady

	server *http.Server
}

func NewRedirectsListener(log *logan.Entry, config RedirectsConfig, investReady InvestReady) *RedirectsListener {
	return &RedirectsListener{
		log:         log.WithField("worker", "redirects_listener"),
		config:      config,
		investReady: investReady,
	}
}

// Run is blocking.
func (l *RedirectsListener) Run(ctx context.Context) {
	l.log.WithField("config", l.config).Info("Starting listening to redirects.")

	go running.UntilSuccess(ctx, l.log, "redirects_listening_server", func(ctx context.Context) (bool, error) {
		l.server = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", l.config.Host, l.config.Port),
			Handler:      http.HandlerFunc(l.redirectsHandler),
			WriteTimeout: l.config.Timeout,
		}

		err := l.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return false, errors.Wrap(err, "Failed to ListenAndServe (Server stopped with error)")
		}

		return false, nil
	}, time.Second, time.Hour)

	<-ctx.Done()

	shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	l.server.Shutdown(shutdownCtx)
	l.log.Info("Server stopped cleanly.")
}

// TODO
func (l *RedirectsListener) redirectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		l.log.WithField("request_method", r.Method).Warn("Received request with wrong method.")
		writeError(w, http.StatusMethodNotAllowed, "Only method POST is allowed.")
		return
	}

	if r.Body == nil {
		l.log.Warn("Received request with empty body.")
		writeError(w, http.StatusBadRequest, "Empty request body.")
		return
	}

	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.log.WithError(err).Warn("Failed to read bytes from request body Reader.")
		writeError(w, http.StatusBadRequest, "Cannot read request body.")
		return
	}
	logger := l.log.WithField("raw_request", string(bb))

	var request redirectedRequest
	err = json.Unmarshal(bb, &request)
	if err != nil {
		logger.WithError(err).Warn("Failed to unmarshal request bytes into struct.")
		writeError(w, http.StatusBadRequest, "Cannot parse JSON request.")
		return
	}

	if validationErr := request.Validate(); validationErr != "" {
		logger.WithField("validation_err", validationErr).Warn("Received invalid request.")
		writeError(w, http.StatusBadRequest, validationErr)
		return
	}

	// TODO

	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, statusCode int, errorMessage string) error {
	resp := struct {
		Error string `json:"error"`
	}{
		Error: errorMessage,
	}

	bb, err := json.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal response to bytes")
	}

	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(bb)
	if err != nil {
		return errors.Wrap(err, "Failed to write marshaled response to the ResponseWriter", logan.F{
			"marshaled_response": string(bb),
		})
	}

	return nil
}
