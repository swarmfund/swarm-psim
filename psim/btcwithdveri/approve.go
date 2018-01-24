package btcwithdveri

import (
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

func (s *Service) preliminaryApproveHandler(w http.ResponseWriter, r *http.Request) {
	approveRequest := withdraw.ApproveRequest{}
	ok := s.readAPIRequest(w, r, &approveRequest)
	if !ok {
		return
	}

	logger := s.log.WithField("preliminary_approve_request", approveRequest)

	// FIXME Change RequestType (to pending)
	checkErr, err := s.obtainAndCheckRequestWithTXHex(approveRequest.Request.ID, approveRequest.Request.Hash, int32(xdr.ReviewableRequestTypeWithdraw), approveRequest.TXHex)
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

	// ApproveRequest is valid
	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.SourceKP).Op(xdrbuild.ReviewRequestOp{
		ID:     approveRequest.Request.ID,
		Hash:   approveRequest.Request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
	}).Sign(s.config.SignerKP).Marshal()
	if err != nil {
		logger.WithError(err).Error("Failed to marshal signed Envelope")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	s.marshalResponseEnvelope(w, r, signedEnvelope)
	logger.Info("Verified Preliminary Approve successfully.")
}

func (s *Service) approveHandler(w http.ResponseWriter, r *http.Request) {
	approveRequest := withdraw.ApproveRequest{}
	ok := s.readAPIRequest(w, r, &approveRequest)
	if !ok {
		return
	}

	logger := s.log.WithField("approve_request", approveRequest)

	checkErr, err := s.obtainAndCheckRequestWithTXHex(approveRequest.Request.ID, approveRequest.Request.Hash, int32(xdr.ReviewableRequestTypeWithdraw), approveRequest.TXHex)
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

	s.marshalResponseEnvelope(w, r, signedEnvelope)
	logger.Info("Verified Approve successfully.")
}

func (s *Service) obtainAndCheckRequestWithTXHex(requestID uint64, requestHash string, neededRequestType int32, btcTXHex string) (checkErr string, err error) {
	addr, amount, checkErr, err := s.obtainAndCheckRequest(requestID, requestHash, neededRequestType)

	validationErr, err := withdraw.ValidateBTCTx(btcTXHex, s.btcClient.GetNetParams(), addr, s.config.HotWalletAddress, amount)
	if err != nil {
		return "", errors.Wrap(err, "Failed to validate BTC TX", logan.F{
			"request_addr":   addr,
			"request_amount": amount,
		})
	}

	return validationErr, nil
}
