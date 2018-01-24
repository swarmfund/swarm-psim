package withdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

func (s *Service) getRejectReason(request horizon.Request) RejectReason {
	address, err := GetWithdrawAddress(request)
	if err != nil {
		switch errors.Cause(err) {
		case ErrMissingAddress:
			return RejectReasonMissingAddress
		case ErrAddressNotAString:
			return RejectReasonAddressNotAString
		}
	}

	err = s.offchainHelper.ValidateAddress(address)
	if err != nil {
		return RejectReasonInvalidAddress
	}

	amount := GetWithdrawAmount(request)
	if amount < s.offchainHelper.GetMinWithdrawAmount() {
		return RejectReasonTooLittleAmount
	}

	return ""
}

func (s *Service) processRequestReject(ctx context.Context, request horizon.Request, reason RejectReason) error {
	returnedEnvelope, err := s.sendRequestToVerify(VerifyRejectURLSuffix, NewReject(request.ID, request.Hash, reason))
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

	s.log.WithFields(GetRequestLoganFields("request", request)).WithField("reject_reason", reason).
		Info("Processed PermanentReject successfully.")

	return nil
}
