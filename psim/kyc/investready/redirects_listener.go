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
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

const (
	UserHashExtDetailsKey = "invest_ready_user_hash"
)

type redirectedRequest struct {
	AccountID string `json:"account_id"`
	OauthCode string `json:"oauth_code"`
}

func (r redirectedRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id": r.AccountID,
	}
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
	log    *logan.Entry
	config RedirectsConfig

	kycRequestsConnector KYCRequestsConnector
	requestPerformer     RequestPerformer
	investReady          InvestReady

	server *http.Server
}

func NewRedirectsListener(
	log *logan.Entry,
	config RedirectsConfig,
	kycRequestsConnector KYCRequestsConnector,
	investReady InvestReady,
	performer RequestPerformer) *RedirectsListener {

	return &RedirectsListener{
		log:                  log.WithField("worker", "redirects_listener"),
		config:               config,
		kycRequestsConnector: kycRequestsConnector,
		requestPerformer:     performer,
		investReady:          investReady,
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

// TODO Tru to refactor me into several methods
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

	var request redirectedRequest
	err = json.Unmarshal(bb, &request)
	if err != nil {
		l.log.WithField("raw_request", string(bb)).WithError(err).Warn("Failed to unmarshal request bytes into struct.")
		writeError(w, http.StatusBadRequest, "Cannot parse JSON request.")
		return
	}
	logger := l.log.WithField("request", request)

	if validationErr := request.Validate(); validationErr != "" {
		l.log.WithField("validation_err", validationErr).Warn("Received invalid request.")
		writeError(w, http.StatusBadRequest, validationErr)
		return
	}

	forbiddenErr, err := l.obtainAndSaveUserHash(r.Context(), request.AccountID, request.OauthCode)
	if err != nil {
		logger.WithError(err).Error("Failed to obtain and save UserHash.")
		writeError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}
	if forbiddenErr != nil {
		logger.WithField("forbidden_reason", forbiddenErr).Warn("User is forbidden to add InvestReady UserHash to the KYCRequest.")
		writeError(w, http.StatusForbidden, "You are not allowed to add InvestReady User to the your KYCRequest.")
		return
	}

	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (l *RedirectsListener) obtainAndSaveUserHash(ctx context.Context, accID, oauthCode string) (forbiddenReason error, err error) {
	kycRequests, err := l.kycRequestsConnector.Requests(fmt.Sprintf("account_to_update_kyc=%s", accID),
		"", horizon.ReviewableRequestType("update_kyc"))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get KYCRequests of the Account from Horizon")
	}
	if len(kycRequests) == 0 {
		return errors.New("No KYCRequests were found for the Account."), nil
	}
	kycRequest := kycRequests[0]
	fields := logan.F{
		"kyc_request": kycRequest,
	}

	if kycRequest.State != kyc.RequestStatePending {
		return errors.Errorf("Expected KYCRequest State to be Pending(%d), but got (%d).", kyc.RequestStatePending, kycRequest.State), nil
	}
	if kycRequest.Details.RequestType != int32(xdr.ReviewableRequestTypeUpdateKyc) {
		return nil, errors.From(errors.Errorf("Expected KYCRequest State to be Pending(%d), but got (%d).", kyc.RequestStatePending, kycRequest.State), fields)
	}
	if kycRequest.Details.KYC.PendingTasks&kyc.TaskUSA == 0 {
		// No job for InvestReady service
		return errors.New("This User's KYCRequest does not have the Task to verify AccreditedInvestor."), nil
	}

	userToken, err := l.investReady.ObtainUserToken(oauthCode)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to obtain User's AccessToken in InvestReady", fields)
	}

	userHash, err := l.investReady.UserHash(userToken)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to ger UserHash in InvestReady", fields)
	}
	fields["user_hash"] = userHash

	err = l.requestPerformer.Approve(ctx, kycRequest.ID, kycRequest.Hash, kyc.TaskCheckInvestReady, kyc.TaskUSA, map[string]string{
		UserHashExtDetailsKey: userHash,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to approve KYCRequest", fields)
	}

	return nil, nil
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
