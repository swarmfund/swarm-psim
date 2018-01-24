package btcwithdveri

import (
	"context"
	"encoding/json"
	"net/http"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

// TODO Pprof
// ServeAPI is blocking method.
func (s *Service) serveAPI(ctx context.Context) {
	r := ape.DefaultRouter()

	r.Post(withdraw.VerifyPreliminaryApproveURLSuffix, s.preliminaryApproveHandler)
	r.Post(withdraw.VerifyApproveURLSuffix, s.approveHandler)
	r.Post(withdraw.VerifyRejectURLSuffix, s.rejectHandler)

	// TODO
	//if s.config.Pprof {
	//	s.log.Info("enabling debugging endpoints")
	//	ape.InjectPprof(r)
	//}

	s.log.WithField("address", s.listener.Addr().String()).Info("Listening.")

	err := ape.ListenAndServe(ctx, s.listener, r)
	if err != nil {
		s.log.WithError(err).Error("ListenAndServe returned error.")
		return
	}
	return
}

func (s *Service) readAPIRequest(w http.ResponseWriter, r *http.Request, request interface{}) (success bool) {
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.log.WithError(err).Warn("Failed to parse request.")
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse JSON request."))
		return false
	}

	return true
}

func (s *Service) obtainAndCheckRequest(requestID uint64, requestHash string, neededRequestType int32) (addr string, amount float64, checkErr string, err error) {
	request, err := withdraw.ObtainRequest(s.horizon.Client(), requestID)
	if err != nil {
		return "", 0, "", errors.Wrap(err, "Failed to Obtain WithdrawRequest from Horizon")
	}

	requestFields := withdraw.GetRequestLoganFields("withdraw_request", *request)

	if request.Hash != requestHash {
		return "", 0, fmt.Sprintf("The RequestHash from Horizon (%s) does not match the one provided (%s).", request.Hash, requestHash), nil
	}
	proveErr := withdraw.ProvePendingRequest(*request, neededRequestType, withdraw.BTCAsset)
	if proveErr != "" {
		return "", 0, fmt.Sprintf("Not a pending BTC WithdrawRequest: %s", proveErr), nil
	}

	addr, err = withdraw.GetWithdrawAddress(*request)
	if err != nil {
		return "", 0, "", errors.Wrap(err, "Failed to get Address from the WithdrawRequest", requestFields)
	}
	amount = withdraw.GetWithdrawAmount(*request)

	return addr, amount, "", nil
}

func (s *Service) marshalResponseEnvelope(w http.ResponseWriter, r *http.Request, envelope string) {
	response := withdraw.EnvelopeResponse{
		Envelope: envelope,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		s.log.WithField("response_trying_to_render", response).WithError(err).Error("Failed to marshal EnvelopeResponse.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	w.Write(respBytes)
}
