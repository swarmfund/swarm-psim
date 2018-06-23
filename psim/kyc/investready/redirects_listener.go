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
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/listener"
)

const (
	UserHashExtDetailsKey = "invest_ready_user_hash"
)

type KYCRequestsConnector interface {
	Requests(filters, cursor string, reqType horizon.ReviewableRequestType) ([]horizon.Request, error)
}

type RedirectsListener struct {
	log    *logan.Entry
	config listener.Config

	kycRequestsConnector KYCRequestsConnector
	requestPerformer     RequestPerformer
	investReady          InvestReady
	doorman              doorman.Doorman

	server *http.Server
}

func NewRedirectsListener(
	log *logan.Entry,
	config listener.Config,
	kycRequestsConnector KYCRequestsConnector,
	investReady InvestReady,
	doorman doorman.Doorman,
	performer RequestPerformer) *RedirectsListener {

	return &RedirectsListener{
		log:                  log.WithField("worker", "redirects_listener"),
		config:               config,
		kycRequestsConnector: kycRequestsConnector,
		requestPerformer:     performer,
		doorman:              doorman,
		investReady:          investReady,
	}
}

// Run is blocking.
func (l *RedirectsListener) Run(ctx context.Context) {
	l.log.WithField("config", l.config).Info("Starting listening to redirects.")

	mux := http.NewServeMux()
	mux.HandleFunc("/", l.redirectsHandler)
	mux.HandleFunc("/user_hash", l.userHashHandler)

	go running.UntilSuccess(ctx, l.log, "redirects_listening_server", func(ctx context.Context) (bool, error) {
		l.server = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", l.config.Host, l.config.Port),
			Handler:      mux,
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

func (l *RedirectsListener) redirectsHandler(w http.ResponseWriter, r *http.Request) {
	bb, errResponseWritten := l.validateHTTPRequest(w, r, http.MethodPost)
	if errResponseWritten {
		return
	}

	var request redirectedRequest
	err := json.Unmarshal(bb, &request)
	if err != nil {
		l.log.WithField("raw_request", string(bb)).WithError(err).Warn("Failed to unmarshal request bytes into struct.")
		writeError(w, http.StatusBadRequest, "Cannot parse JSON request.")
		return
	}

	l.processRedirectRequest(r.Context(), w, request)
}

func (l *RedirectsListener) processRedirectRequest(ctx context.Context, w http.ResponseWriter, request redirectedRequest) {
	logger := l.log.WithField("request", request)

	if validationErr := request.Validate(); validationErr != "" {
		logger.WithField("validation_err", validationErr).Warn("Received invalid request.")
		writeError(w, http.StatusBadRequest, validationErr)
		return
	}

	kycRequest, forbiddenErr, err := l.getKYCRequest(ctx, request.AccountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get KYCRequest by AccountID.")
		writeError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}
	if forbiddenErr != nil {
		logger.WithField("forbidden_reason", forbiddenErr).Warn("User is forbidden to add InvestReady UserHash to the KYCRequest.")
		writeError(w, http.StatusForbidden, forbiddenErr.Error())
		return
	}

	err = l.obtainAndSaveUserHash(ctx, *kycRequest, request.AccountID, request.OauthCode)
	if err != nil {
		logger.WithError(err).Error("Failed to obtain and save UserHash.")
		writeError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}

	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (l *RedirectsListener) getKYCRequest(ctx context.Context, accID string) (kycRequest *horizon.Request, forbiddenReason error, err error) {
	kycRequests, err := l.kycRequestsConnector.Requests(fmt.Sprintf("account_to_update_kyc=%s", accID),
		"", horizon.ReviewableRequestType("update_kyc"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to get KYCRequests of the Account from Horizon")
	}
	if len(kycRequests) == 0 {
		return nil, errors.New("No KYCRequests were found for the Account."), nil
	}
	r := kycRequests[0]
	kycRequest = &r
	fields := logan.F{
		"kyc_request": kycRequest,
	}

	if kycRequest.State != kyc.RequestStatePending {
		return nil, errors.Errorf("Expected KYCRequest State to be Pending(%d), but got (%d).", kyc.RequestStatePending, kycRequest.State), nil
	}
	if kycRequest.Details.RequestType != int32(xdr.ReviewableRequestTypeUpdateKyc) {
		return nil, nil, errors.From(errors.Errorf("Expected Request type to be KYC(%d), but got (%d).",
			xdr.ReviewableRequestTypeUpdateKyc, kycRequest.Details.RequestType), fields)
	}
	if kycRequest.Details.KYC.PendingTasks&kyc.TaskUSA == 0 {
		// No job for InvestReady service
		return nil, errors.New("This User's KYCRequest does not have the Task to verify AccreditedInvestor."), nil
	}
	if kycRequest.Details.KYC.PendingTasks&(kyc.TaskSuperAdmin|kyc.TaskFaceValidation|kyc.TaskDocsExpirationDate|kyc.TaskSubmitIDMind|kyc.TaskCheckIDMind) != 0 {
		// Some previous jobs are not finished
		return nil, errors.New("This User's KYCRequest has some other Tasks to be done first."), nil
	}

	return kycRequest, nil, nil
}

func (l *RedirectsListener) obtainAndSaveUserHash(ctx context.Context, kycRequest horizon.Request, accID, oauthCode string) error {
	userToken, err := l.investReady.ObtainUserToken(oauthCode)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain User's AccessToken in InvestReady")
	}

	userHash, err := l.investReady.UserHash(userToken)
	if err != nil {
		return errors.Wrap(err, "Failed to ger UserHash in InvestReady")
	}
	fields := logan.F{
		"user_hash": userHash,
	}

	err = l.saveUserHash(ctx, kycRequest, accID, userHash)
	if err != nil {
		return errors.Wrap(err, "Failed to approve KYCRequest (with saving UserHash)", fields)
	}

	return nil
}

func (l *RedirectsListener) saveUserHash(ctx context.Context, kycRequest horizon.Request, accID, userHash string) error {
	err := l.requestPerformer.Approve(ctx, kycRequest.ID, kycRequest.Hash, kyc.TaskCheckInvestReady, kyc.TaskUSA, map[string]string{
		UserHashExtDetailsKey: userHash,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to approve KYCRequest")
	}

	l.log.WithFields(logan.F{
		"account_id": accID,
		"user_hash":  userHash,
	}).Info("Saved UserHash into Core successfully.")
	return nil
}

func (l *RedirectsListener) validateHTTPRequest(w http.ResponseWriter, r *http.Request, requestMethod string) (respBody []byte, errResponseWritten bool) {
	if r.Method != requestMethod {
		l.log.WithField("request_method", r.Method).Warn("Received request with wrong method.")
		writeError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Only method %s is allowed.", requestMethod))
		return nil, true
	}

	if r.Body == nil {
		l.log.Warn("Received request with empty body.")
		writeError(w, http.StatusBadRequest, "Empty request body.")
		return nil, true
	}

	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.log.WithError(err).Warn("Failed to read bytes from request body Reader.")
		writeError(w, http.StatusBadRequest, "Cannot read request body.")
		return nil, true
	}

	if l.config.CheckSignature {
		var request struct {
			AccountID string `json:"account_id"`
		}
		err = json.Unmarshal(bb, &request)
		if err != nil {
			l.log.WithField("raw_request", string(bb)).WithError(err).Warn("Failed to preliminary unmarshal request bytes into struct(with only AccountID).")
			writeError(w, http.StatusBadRequest, "Cannot parse JSON request.")
			return nil, true
		}

		err := l.doorman.Check(r, doorman.SignatureOf(request.AccountID))
		if err != nil {
			l.log.WithError(err).Warn("Request signature is invalid.")
			writeError(w, http.StatusUnauthorized, err.Error())
			return nil, true
		}
	}

	return bb, false
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
