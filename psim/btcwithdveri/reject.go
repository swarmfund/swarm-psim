package btcwithdveri

import (
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

func (s *Service) rejectHandler(w http.ResponseWriter, r *http.Request) {
	rejectRequest := withdraw.RejectRequest{}
	ok := s.readAPIRequest(w, r, &rejectRequest)
	if !ok {
		return
	}

	logger := s.log.WithField("reject_request", rejectRequest)

	addr, amount, checkErr, err := s.obtainAndCheckRequest(rejectRequest.Request.ID, rejectRequest.Request.Hash, int32(xdr.ReviewableRequestTypeWithdraw))
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

	logger = logger.WithFields(logan.F{
		"withdraw_address": addr,
		"withdraw_amount":  amount,
	})

	validationErr := s.validateReject(addr, amount, rejectRequest.RejectReason)
	if validationErr != "" {
		logger.WithField("validation_error", validationErr).Warn("Got invalid RejectReason.")
		ape.RenderErr(w, r, problems.Forbidden(validationErr))
		return
	}

	// RejectRequest is valid
	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.SourceKP).Op(xdrbuild.ReviewRequestOp{
		ID:     rejectRequest.Request.ID,
		Hash:   rejectRequest.Request.Hash,
		Action: xdr.ReviewRequestOpActionPermanentReject,
		Details: xdrbuild.WithdrawalDetails{
			ExternalDetails: "",
		},
		Reason: string(rejectRequest.RejectReason),
	}).Sign(s.config.SignerKP).Marshal()
	if err != nil {
		logger.WithError(err).Error("Failed to marshal signed Envelope.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	s.marshalResponseEnvelope(w, r, signedEnvelope)
	logger.Info("Verified Reject successfully.")
}

func (s *Service) validateReject(addr string, amount float64, reason withdraw.RejectReason) string {
	switch reason {
	case withdraw.RejectReasonInvalidAddress:
		return s.validateInvalidAddress(addr)
	case withdraw.RejectReasonTooLittleAmount:
		return s.validateTooLittleAmount(amount)
	}

	return fmt.Sprintf("Unsupported RejectReason '%s'.", string(reason))
}

func (s *Service) validateInvalidAddress(addr string) string {
	err := withdraw.ValidateBTCAddress(addr, s.btcClient.GetNetParams())
	if err != nil {
		return ""
	} else {
		return fmt.Sprintf("BTC Address '%s' is actually valid.", addr)
	}
}

func (s *Service) validateTooLittleAmount(amount float64) string {
	if amount < s.config.MinWithdrawAmount {
		return ""
	} else {
		return fmt.Sprintf("Amount is not actually little. I consider %f as MinWithdrawAmount.", s.config.MinWithdrawAmount)
	}
}
