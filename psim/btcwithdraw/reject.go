package btcwithdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

func (s *Service) getRejectReason(withdrawAddress string, amount float64) withdraw.RejectReason {
	err := withdraw.ValidateBTCAddress(withdrawAddress, s.btcClient.GetNetParams())
	if err != nil {
		return withdraw.RejectReasonInvalidAddress
	}

	if amount < s.config.MinWithdrawAmount {
		return withdraw.RejectReasonTooLittleAmount
	}

	return ""
}

func (s *Service) verifyReject(ctx context.Context, request horizon.Request, reason withdraw.RejectReason) error {
	returnedEnvelope, err := s.sendRequestToVerify(withdraw.VerifyRejectURLSuffix, withdraw.NewReject(request.ID, request.Hash, reason))
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

	s.log.WithFields(withdraw.GetRequestLoganFields("request", request)).WithField("reject_reason", reason).
		Info("Sent PermanentReject to Verify successfully.")
	return nil
}
