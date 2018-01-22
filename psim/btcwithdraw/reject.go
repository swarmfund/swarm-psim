package btcwithdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

func (s *Service) getRejectReason(withdrawAddress string, amount float64) RejectReason {
	err := ValidateBTCAddress(withdrawAddress, s.btcClient.GetNetParams())
	if err != nil {
		return RejectReasonInvalidAddress
	}

	if amount < s.config.MinWithdrawAmount {
		return RejectReasonTooLittleAmount
	}

	return ""
}

func (s *Service) verifyReject(ctx context.Context, request horizon.Request, reason RejectReason) error {
	returnedEnvelope, err := s.sendRequestToVerify(RejectRequest{
		Request: WithdrawalRequest{
			ID:   request.ID,
			Hash: request.Hash,
		},
		RejectReason: reason,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to send Reject to Verify")
	}

	checkErr := checkRejectEnvelope(*returnedEnvelope, request.ID, request.Hash, reason)
	if checkErr != "" {
		return errors.Wrap(err, "Envelope returned by Verify is invalid")
	}

	s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope")
	}

	s.log.WithFields(GetRequestLoganFields("request", request)).WithField("reject_reason", reason).
		Info("Sent PermanentReject to Verify successfully.")
	return nil
}
