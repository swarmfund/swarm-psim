package withdveri

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

	checkErr, err := s.obtainAndCheckRequestWithTXHex(approveRequest.Request.ID, approveRequest.Request.Hash, int32(xdr.ReviewableRequestTypeTwoStepWithdrawal), approveRequest.TXHex)
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
	extDetails := withdraw.ExternalDetails{
		TXHex: approveRequest.TXHex,
	}
	extDetBytes, err := json.Marshal(extDetails)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal ExternalDetails into JSON.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	signedEnvelope, err := s.xdrbuilder.Transaction(s.sourceKP).Op(xdrbuild.ReviewRequestOp{
		ID:     approveRequest.Request.ID,
		Hash:   approveRequest.Request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.TwoStepWithdrawalDetails{
			ExternalDetails: string(extDetBytes),
		},
	}).Sign(s.signerKP).Marshal()
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

	fullySignedOffchainTXHex, err := s.offchainHelper.SignTX(approveRequest.TXHex)
	if err != nil {
		logger.WithError(err).Error("Failed to sign the TX.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	extDetails := withdraw.ExternalDetails{
		TXHex: fullySignedOffchainTXHex,
	}
	extDetBytes, err := json.Marshal(extDetails)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal ExternalDetails into JSON.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	// ApproveRequest is valid
	signedEnvelope, err := s.xdrbuilder.Transaction(s.sourceKP).Op(xdrbuild.ReviewRequestOp{
		ID:     approveRequest.Request.ID,
		Hash:   approveRequest.Request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.WithdrawalDetails{
			ExternalDetails: string(extDetBytes),
		},
	}).Sign(s.signerKP).Marshal()
	if err != nil {
		logger.WithError(err).Error("Failed to marshal signed Envelope.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	s.marshalResponseEnvelope(w, r, signedEnvelope)
	logger.Info("Verified Approve successfully.")
}

func (s *Service) obtainAndCheckRequestWithTXHex(requestID uint64, requestHash string, neededRequestType int32, txHex string) (checkErr string, err error) {
	request, checkErr, err := s.obtainAndCheckRequest(requestID, requestHash, neededRequestType)
	if err != nil {
		return "", errors.Wrap(err, "Failed to obtain-and-check Request")
	}
	if checkErr != "" {
		return checkErr, nil
	}

	addr, err := withdraw.GetWithdrawAddress(*request)
	if err != nil {
		return err.Error(), nil
	}
	amount, err := withdraw.GetWithdrawAmount(*request)
	if err != nil {
		return err.Error(), nil
	}
	amount = s.offchainHelper.ConvertAmount(amount)

	validationErr, err := s.offchainHelper.ValidateTX(txHex, addr, amount)
	if err != nil {
		return "", errors.Wrap(err, "Failed to validate the TX", logan.F{
			"request_addr":   addr,
			"request_amount": amount,
		})
	}

	return validationErr, nil
}
