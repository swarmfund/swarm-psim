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
	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/swarmfund/go/amount"
)

func (s *Service) processReject(w http.ResponseWriter, r *http.Request, withdrawRequest *horizonV2.Request, horizonTX *horizon.TransactionBuilder) {
	opBody := horizonTX.Operations[0].Body.ReviewRequestOp
	fields := logan.F{
		"request_id": opBody.RequestId,
		"request_hash": opBody.RequestHash,
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

func (s *Service) validateReject(withdrawRequest *horizonV2.Request, reason btcwithdraw.RejectReason) string {
	switch reason {
	case btcwithdraw.RejectReasonInvalidAddress:
		addr := string(withdrawRequest.Details.Withdraw.ExternalDetails)
		_, err := btc.NewAddrFromString(addr)
		if err != nil {
			return ""
		} else {
			return fmt.Sprintf("BTC Address '%s' is actually valid.", addr)
		}
	case btcwithdraw.RejectReasonTooLittleAmount:
		withdrawAmount := float64(int64(withdrawRequest.Details.Withdraw.DestinationAmount)) / amount.One
		if withdrawAmount < s.config.MinWithdrawAmount {
			return ""
		} else {
			return fmt.Sprintf("Amount is not actually little. I consider %f as MinWithdrawAmount.", s.config.MinWithdrawAmount)
		}
	}

	return fmt.Sprintf("Unsupported RejectReason '%s'.", string(reason))
}

// FIXME
func (s *Service) processValidReject(tx *horizon.TransactionBuilder) error {
	// FIXME
	return tx.Sign(s.config.SignerKP).Submit()
}
