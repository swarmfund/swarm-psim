package withdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources"
)

func (s *Service) getRejectReason(request regources.ReviewableRequest) RejectReason {
	address, err := GetWithdrawalAddress(request)
	if err != nil {
		switch errors.Cause(err) {
		case ErrMissingTwoStepWithdraw, ErrMissingAddress:
			return RejectReasonMissingAddress
		case ErrAddressNotAString:
			return RejectReasonAddressNotAString
		}
	}

	err = s.offchainHelper.ValidateAddress(address)
	if err != nil {
		return RejectReasonInvalidAddress
	}

	destAmount, err := GetWithdrawAmount(request)
	if err != nil {
		return RejectReasonMissingAmount
	}

	amount := s.offchainHelper.ConvertAmount(destAmount)
	if amount < s.offchainHelper.GetMinWithdrawAmount() {
		return RejectReasonTooLittleAmount
	}

	return ""
}

func (s *Service) processRequestReject(ctx context.Context, request regources.ReviewableRequest, reason RejectReason) error {
	returnedEnvelope, err := s.sendRequestToVerifier(VerifyRejectURLSuffix, NewReject(request.ID, request.Hash, reason))
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

	err = s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope", logan.F{
			"envelope_returned_by_verify": returnedEnvelope,
		})
	}

	s.log.WithField("request", request).WithField("reject_reason", reason).
		Info("Processed PermanentReject successfully.")

	return nil
}
