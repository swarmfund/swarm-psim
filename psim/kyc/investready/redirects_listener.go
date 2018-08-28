package investready

import (
	"context"
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/swarmfund/psim/psim/listener"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

const (
	UserHashExtDetailsKey = "invest_ready_user_hash"
)

type KYCRequestsConnector interface {
	Requests(filters, cursor string, reqType horizon.ReviewableRequestType) ([]regources.ReviewableRequest, error)
	GetRequestByID(requestID uint64) (*regources.ReviewableRequest, error)
}

type AccountsConnector interface {
	ByAddress(address string) (*horizon.Account, error)
}

type RedirectsListener struct {
	log    *logan.Entry
	config listener.Config

	kycRequestsConnector KYCRequestsConnector
	accountsConnector    AccountsConnector
	requestPerformer     RequestPerformer
	investReady          InvestReady
	doorman              doorman.Doorman

	server *http.Server
}

func NewRedirectsListener(
	log *logan.Entry,
	config listener.Config,
	kycRequestsConnector KYCRequestsConnector,
	accountsConnector AccountsConnector,
	investReady InvestReady,
	doorman doorman.Doorman,
	performer RequestPerformer) *RedirectsListener {

	return &RedirectsListener{
		log:                  log.WithField("worker", "redirects_listener"),
		config:               config,
		kycRequestsConnector: kycRequestsConnector,
		accountsConnector:    accountsConnector,
		requestPerformer:     performer,
		doorman:              doorman,
		investReady:          investReady,
	}
}

// Run is blocking.
func (l *RedirectsListener) Run(ctx context.Context) {
	l.log.WithField("config", l.config).Info("Starting listening to redirects.")

	r := chi.NewRouter()
	r.Put("/user_hash", l.userHashHandler)
	r.Post("/*", l.redirectsHandler)

	listener.RunServer(ctx, l.log, r, l.config)
}

func (l *RedirectsListener) redirectsHandler(w http.ResponseWriter, r *http.Request) {
	bb, errResponseWritten := listener.ValidateHTTPRequest(w, r, l.log, l.doorman)
	if errResponseWritten {
		return
	}

	var request redirectedRequest
	err := json.Unmarshal(bb, &request)
	if err != nil {
		l.log.WithField("raw_request", string(bb)).WithError(err).Warn("Failed to unmarshal request bytes into struct.")
		listener.WriteError(w, http.StatusBadRequest, "Cannot parse JSON request.")
		return
	}

	l.processRedirectRequest(r.Context(), w, request)
}

func (l *RedirectsListener) processRedirectRequest(ctx context.Context, w http.ResponseWriter, req redirectedRequest) {
	logger := l.log.WithField("request", req)

	if validationErr := req.Validate(); validationErr != "" {
		logger.WithField("validation_err", validationErr).Warn("Received invalid request.")
		listener.WriteError(w, http.StatusBadRequest, validationErr)
		return
	}

	kycRequest, forbiddenErr, err := l.getAndValidateKYCRequest(ctx, req.AccountID, req.KYCRequestID)
	if err != nil {
		logger.WithError(err).Error("Failed to get KYCRequest by AccountID.")
		listener.WriteError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}
	if forbiddenErr != nil {
		logger.WithField("forbidden_reason", forbiddenErr).Warn("User is forbidden to add InvestReady UserHash to the KYCRequest.")
		listener.WriteError(w, http.StatusForbidden, forbiddenErr.Error())
		return
	}

	err = l.obtainAndSaveUserHash(ctx, req.OauthCode, *kycRequest, req.AccountID)
	if err != nil {
		logger.WithError(err).Error("Failed to obtain and save UserHash.")
		listener.WriteError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}

	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (l *RedirectsListener) getAndValidateKYCRequest(ctx context.Context, accID string, kycRequestID uint64) (kycRequest *regources.ReviewableRequest, forbiddenReason error, err error) {
	if vErr, err := l.validateAccount(accID); err != nil || vErr != nil {
		return nil, vErr, errors.Wrap(err, "failed to validate Account")
	}

	kycRequest, err = l.kycRequestsConnector.GetRequestByID(kycRequestID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to get KYCRequest by ID")
	}
	fields := logan.F{
		"kyc_request": kycRequest,
	}

	if err := l.validateKYCRequest(*kycRequest); err != nil {
		return nil, errors.From(err, fields), nil
	}

	return kycRequest, nil, nil
}

func (l *RedirectsListener) validateAccount(accID string) (forbiddenReason error, err error) {
	account, err := l.accountsConnector.ByAddress(accID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}
	if account.IsBlocked {
		l.log.WithField("account_id", accID).Debug("skipping, since Account is blocked")
		return errors.New("account is blocked"), nil
	}
	if !kyc.IsVerified(*account) {
		return errors.New("only Accounts of Type Verified are allowed to upgrade to General via InvestReady"), nil
	}

	return nil, nil
}

func (l *RedirectsListener) validateKYCRequest(kycRequest regources.ReviewableRequest) (validationErr error) {
	if kycRequest.State != kyc.RequestStatePending {
		return errors.Errorf("Expected KYCRequest State to be Pending(%d), but got (%d).", kyc.RequestStatePending, kycRequest.State)
	}
	if kycRequest.Details.RequestType != int32(xdr.ReviewableRequestTypeUpdateKyc) {
		return errors.Errorf("Expected Request type to be KYC(%d), but got (%d).",
			xdr.ReviewableRequestTypeUpdateKyc, kycRequest.Details.RequestType)
	}
	if !kyc.IsUpdateToGeneral(kycRequest) {
		return errors.New("AccountTypeToSet must be (General)")
	}
	if kycRequest.Details.KYC.PendingTasks&kyc.TaskUSA == 0 {
		// No job for InvestReady service
		return errors.New("This User's KYCRequest does not have the Task to verify AccreditedInvestor.")
	}
	if kycRequest.Details.KYC.PendingTasks&(kyc.TaskSuperAdmin|kyc.TaskFaceValidation|kyc.TaskDocsExpirationDate|kyc.TaskSubmitIDMind|kyc.TaskCheckIDMind) != 0 {
		// Some previous jobs are not finished
		return errors.New("This User's KYCRequest has some other Tasks to be done first.")
	}

	return nil
}

func (l *RedirectsListener) obtainAndSaveUserHash(ctx context.Context, oauthCode string, kycRequest regources.ReviewableRequest, accID string) error {
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

	err = l.approveRequestAddUserHash(ctx, kycRequest, accID, userHash)
	if err != nil {
		return errors.Wrap(err, "Failed to approve KYCRequest (with saving UserHash)", fields)
	}

	return nil
}

func (l *RedirectsListener) approveRequestAddUserHash(ctx context.Context, kycRequest regources.ReviewableRequest, accID, userHash string) error {
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
