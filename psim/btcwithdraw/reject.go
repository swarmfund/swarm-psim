package btcwithdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
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

func (s *Service) processRequestReject(ctx context.Context, request horizon.Request, reason withdraw.RejectReason) error {
	returnedEnvelope, err := s.sendRequestToVerify(withdraw.VerifyRejectURLSuffix, withdraw.NewReject(request.ID, request.Hash, reason))
	if err != nil {
		return errors.Wrap(err, "Failed to send Reject to Verify")
	}

	checkErr := checkRejectEnvelope(*returnedEnvelope, request.ID, request.Hash, reason)
	if checkErr != "" {
		return errors.From(errors.New("Envelope returned by Verify is invalid."), logan.F{
			"check_error":                 checkErr,
			"envelope_returned_by_verify": returnedEnvelope,
		})
	}

	s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope", logan.F{
			"envelope_returned_by_verify": returnedEnvelope,
		})
	}

	s.log.WithFields(withdraw.GetRequestLoganFields("request", request)).WithField("reject_reason", reason).
		Info("Processed PermanentReject successfully.")

	return nil
}
