package btcwithdveri

import (
	"context"
	"encoding/json"
	"net/http"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
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
	//r.Post(withdraw.VerifyRejectURLSuffix, s.rejectHandler)

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

func (s *Service) preliminaryApproveHandler(w http.ResponseWriter, r *http.Request) {
	approveRequest := withdraw.ApproveRequest{}
	ok := s.readRequest(w, r, &approveRequest)
	if !ok {
		return
	}

	// FIXME Change RequestType (to pending)
	checkErr, err := s.checkWithdrawRequest(approveRequest.Request.ID, approveRequest.Request.Hash, int32(xdr.ReviewableRequestTypeWithdraw), approveRequest.TXHex)
	if err != nil {
		s.log.WithField("preliminary_approve_request", approveRequest).WithError(err).Error("Failed to check WithdrawRequest.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	if checkErr != "" {
		s.log.WithField("preliminary_approve_request", approveRequest).WithField("check_error", checkErr).Warn("Got invalid PreliminaryApproveRequest.")
		ape.RenderErr(w, r, problems.Forbidden(checkErr))
		return
	}

	// ApproveRequest is valid
	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.SourceKP).Op(xdrbuild.ReviewRequestOp{
		ID:     approveRequest.Request.ID,
		Hash:   approveRequest.Request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
	}).Sign(s.config.SignerKP).Marshal()
	if err != nil {
		s.log.WithField("preliminary_approve_request", approveRequest).WithError(err).Error("Failed to marshal signed Envelope")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	response := withdraw.EnvelopeResponse{
		Envelope: signedEnvelope,
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		s.log.WithField("response_trying_to_render", response).WithError(err).Error("Failed to marshal response.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	w.Write(respBytes)
}

func (s *Service) approveHandler(w http.ResponseWriter, r *http.Request) {
	approveRequest := withdraw.ApproveRequest{}
	ok := s.readRequest(w, r, &approveRequest)
	if !ok {
		return
	}

	logger := s.log.WithField("approve_request", approveRequest)

	checkErr, err := s.checkWithdrawRequest(approveRequest.Request.ID, approveRequest.Request.Hash, int32(xdr.ReviewableRequestTypeWithdraw), approveRequest.TXHex)
	if err != nil {
		logger.WithError(err).Error("Failed to check WithdrawRequest.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	if checkErr != "" {
		logger.WithField("check_error", checkErr).Warn("Got invalid PreliminaryApproveRequest.")
		ape.RenderErr(w, r, problems.Forbidden(checkErr))
		return
	}

	fullySignedBtcTXHex, err := s.btcClient.SignAllTXInputs(approveRequest.TXHex, s.config.HotWalletScriptPubKey, s.config.HotWalletRedeemScript, s.config.PrivateKey)
	if err != nil {
		logger.WithError(err).Error("Failed to sign BTC TX.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	extDetails := withdraw.ExternalDetails{
		TXHex: fullySignedBtcTXHex,
	}
	extDetBytes, err := json.Marshal(extDetails)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal ExternalDetails into JSON.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	// ApproveRequest is valid
	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.SourceKP).Op(xdrbuild.ReviewRequestOp{
		ID:     approveRequest.Request.ID,
		Hash:   approveRequest.Request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.WithdrawalDetails{
			ExternalDetails: string(extDetBytes),
		},
	}).Sign(s.config.SignerKP).Marshal()
	if err != nil {
		logger.WithError(err).Error("Failed to marshal signed Envelope.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	response := withdraw.EnvelopeResponse{
		Envelope: signedEnvelope,
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		s.log.WithField("response_trying_to_render", response).WithError(err).Error("Failed to marshal response.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	w.Write(respBytes)
}

func (s *Service) readRequest(w http.ResponseWriter, r *http.Request, request interface{}) (success bool) {
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.log.WithError(err).Warn("Failed to parse request.")
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse JSON request."))
		return false
	}

	return true
}

func (s *Service) checkWithdrawRequest(requestID uint64, requestHash string, neededRequestType int32, btcTXHex string) (checkErr string, err error) {
	request, err := withdraw.ObtainRequest(s.horizon.Client(), requestID)
	if err != nil {
		return "", errors.Wrap(err, "Failed to Obtain WithdrawRequest from Horizon")
	}

	requestFields := withdraw.GetRequestLoganFields("withdraw_request", *request)

	if request.Hash != requestHash {
		return fmt.Sprintf("The RequestHash from Horizon (%s) does not match from the one provided (%s).", request.Hash, requestHash), nil
	}
	proveErr := withdraw.ProvePendingBTCRequest(*request, neededRequestType)
	if proveErr != "" {
		return proveErr, nil
	}

	addr, err := withdraw.GetWithdrawAddress(*request)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Address from the WithdrawRequest", requestFields)
	}
	amount := withdraw.GetWithdrawAmount(*request)

	validationErr, err := withdraw.ValidateBTCTx(btcTXHex, s.btcClient.GetNetParams(), addr, s.config.HotWalletAddress, amount)
	if err != nil {
		return "", errors.Wrap(err, "Failed to validate BTC TX", requestFields)
	}

	return validationErr, nil
}
