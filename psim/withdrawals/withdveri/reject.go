package withdveri

import (
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/verification"
	"gitlab.com/swarmfund/psim/psim/withdrawals/withdraw"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
)

func (s *Service) rejectHandler(w http.ResponseWriter, r *http.Request) {
	rejectRequest := withdraw.RejectRequest{}
	if ok := verification.ReadAPIRequest(s.log, w, r, &rejectRequest); !ok {
		return
	}

	logger := s.log.WithField("reject_request", rejectRequest)

	request, checkErr, err := s.obtainAndCheckRequest(rejectRequest.Request.ID, rejectRequest.Request.Hash, int32(xdr.ReviewableRequestTypeTwoStepWithdrawal))
	if err != nil {
		logger.WithError(err).Error("Failed to obtain-and-check WithdrawRequest.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	if checkErr != "" {
		logger.WithField("check_error", checkErr).Warn("Got invalid RejectRequest.")
		ape.RenderErr(w, r, problems.Forbidden(checkErr))
		return
	}

	logger = logger.WithField("request", *request)

	validationErr := s.validateRejectReason(*request, rejectRequest.RejectReason)
	if validationErr != "" {
		logger.WithField("validation_error", validationErr).Warn("Got invalid RejectReason.")
		ape.RenderErr(w, r, problems.Forbidden(validationErr))
		return
	}

	// RejectRequest is valid
	signedEnvelope, err := s.xdrbuilder.Transaction(s.sourceKP).Op(xdrbuild.ReviewRequestOp{
		ID:     rejectRequest.Request.ID,
		Hash:   rejectRequest.Request.Hash,
		Action: xdr.ReviewRequestOpActionPermanentReject,
		Details: xdrbuild.TwoStepWithdrawalDetails{
			ExternalDetails: "",
		},
		Reason: string(rejectRequest.RejectReason),
	}).Sign(s.signerKP).Marshal()
	if err != nil {
		logger.WithError(err).Error("Failed to marshal signed Envelope.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	ok := verification.RenderResponseEnvelope(logger, w, r, signedEnvelope)
	if ok {
		logger.Info("Verified Reject successfully.")
	}
}

func (s *Service) validateRejectReason(request horizon.Request, reason withdraw.RejectReason) string {
	var gettingAddressErr withdraw.RejectReason
	addr, err := withdraw.GetWithdrawalAddress(request)
	if err != nil {
		switch errors.Cause(err) {
		case withdraw.ErrMissingTwoStepWithdraw, withdraw.ErrMissingAddress:
			gettingAddressErr = withdraw.RejectReasonMissingAddress
		case withdraw.ErrAddressNotAString:
			gettingAddressErr = withdraw.RejectReasonAddressNotAString
		}
	}

	if gettingAddressErr != "" {
		if gettingAddressErr == reason {
			return ""
		} else {
			return fmt.Sprintf("Expected RejectReason to be (%s), but received (%s)", gettingAddressErr, reason)
		}
	}

	// No problems getting Address

	switch reason {
	case withdraw.RejectReasonMissingAddress:
		return "Address field is actually present in the details."
	case withdraw.RejectReasonAddressNotAString:
		return "Address from details is actually a string."
	case withdraw.RejectReasonInvalidAddress:
		return s.validateInvalidAddress(addr)
	case withdraw.RejectReasonTooLittleAmount:
		amount := s.offchainHelper.ConvertAmount(int64(request.Details.TwoStepWithdraw.DestAssetAmount))
		return s.validateTooLittleAmount(amount)
	}

	return fmt.Sprintf("Unsupported RejectReason '%s'.", string(reason))
}

func (s *Service) validateInvalidAddress(addr string) string {
	err := s.offchainHelper.ValidateAddress(addr)
	if err != nil {
		return ""
	} else {
		return fmt.Sprintf("The Address '%s' is actually valid.", addr)
	}
}

func (s *Service) validateTooLittleAmount(amount int64) string {
	if amount < s.offchainHelper.GetMinWithdrawAmount() {
		return ""
	} else {
		return fmt.Sprintf("Amount (%d) is not actually little. I consider (%d) as MinWithdrawAmount.", amount, s.offchainHelper.GetMinWithdrawAmount())
	}
}
