package btcwithdveri

import (
	"gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/btcwithdraw"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"net/http"
	"gitlab.com/distributed_lab/logan/v3"
	"fmt"
	"gitlab.com/swarmfund/go/amount"
)

func (s *Service) processReject(w http.ResponseWriter, r *http.Request, withdrawRequest horizonV2.Request, horizonTX *horizon.TransactionBuilder) {
	opBody := horizonTX.Operations[0].Body.ReviewRequestOp
	fields := logan.F{
		"request_id": opBody.RequestId,
		"request_hash": string(opBody.RequestHash[:]),
		"request_action_i": int32(opBody.Action),
		"request_action": opBody.Action.String(),
		"reject_reason": opBody.Reason,
	}

	validationErr := s.validateReject(withdrawRequest, btcwithdraw.RejectReason(opBody.Reason))
	if validationErr != "" {
		ape.RenderErr(w, r, problems.Forbidden(validationErr))
		return
	}

	err := s.processValidReject(horizonTX)
	if err != nil {
		s.log.WithFields(fields).WithError(err).Error("Failed to process valid PermanentReject Request.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	s.log.WithFields(fields).Info("Verified PermanentReject for WithdrawalRequest successfully.")
}

func (s *Service) validateReject(request horizonV2.Request, reason btcwithdraw.RejectReason) string {
	switch reason {
	case btcwithdraw.RejectReasonInvalidAddress:
		return s.validateInvalidAddress(request)
	case btcwithdraw.RejectReasonTooLittleAmount:
		return s.validateTooLittleAmount(request)
	}

	return fmt.Sprintf("Unsupported RejectReason '%s'.", string(reason))
}

func (s *Service) validateInvalidAddress(request horizonV2.Request) string {
	addr, err := btcwithdraw.GetWithdrawAddress(request)
	if err != nil {
		return "Unable to obtain BTC Address of the Withdrawal"
	}

	err = btcwithdraw.ValidateBTCAddress(addr, s.btcClient.GetNetParams())
	if err != nil {
		return ""
	} else {
		return fmt.Sprintf("BTC Address '%s' is actually valid.", addr)
	}
}

func (s *Service) validateTooLittleAmount(request horizonV2.Request) string {
	withdrawAmount := float64(int64(request.Details.Withdraw.DestinationAmount)) / amount.One

	if withdrawAmount < s.config.MinWithdrawAmount {
		return ""
	} else {
		return fmt.Sprintf("Amount is not actually little. I consider %f as MinWithdrawAmount.", s.config.MinWithdrawAmount)
	}
}

func (s *Service) processValidReject(tx *horizon.TransactionBuilder) error {
	return tx.Sign(s.config.SignerKP).Submit()
}
