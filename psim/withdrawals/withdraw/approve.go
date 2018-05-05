package withdraw

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

// ProcessValidPendingRequest knows how to process both TwoStepWithdrawal and Withdraw RequestTypes.
func (s *Service) processValidPendingRequest(ctx context.Context, request horizon.Request) error {
	withdrawAddress, err := GetWithdrawAddress(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get Withdraw Address")
	}
	amount, err := GetWithdrawAmount(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get Withdraw Amount")
	}
	withdrawAmount := s.offchainHelper.ConvertAmount(amount)

	if request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		// TwoStepWithdrawal needs PreliminaryApprove first.
		unsignedOffchainTXHex, err := s.offchainHelper.CreateTX(ctx, withdrawAddress, withdrawAmount)
		if err != nil {
			return errors.Wrap(err, "Failed to create Offchain TX")
		}

		err = s.processPreliminaryApprove(ctx, request, unsignedOffchainTXHex)
		if err != nil {
			return errors.Wrap(err, "Failed to verify Preliminary Approve Request", logan.F{"tx_hex": unsignedOffchainTXHex})
		}
	}

	newRequest, err := s.requestsConnector.GetRequestByID(request.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain Request from Horizon")
	}

	for newRequest.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		s.log.WithField("new_request_id", newRequest.ID).
			Debugf("WithdrawRequest still hasn't changed type to Withdraw(%d). Sleeping for 3 seconds.", xdr.ReviewableRequestTypeWithdraw)
		// TODO Incremental
		time.Sleep(3 * time.Second)
		newRequest, err = s.requestsConnector.GetRequestByID(request.ID)
		if err != nil {
			return errors.Wrap(err, "Failed to obtain Request from Horizon")
		}
	}

	unsignedOffchainTX, err := GetTXHex(*newRequest)
	if err != nil {
		return errors.Wrap(err, "Failed to get TX hex from the WithdrawRequest")
	}

	partlySignedOffchainTX, err := s.offchainHelper.SignTX(unsignedOffchainTX)
	if err != nil {
		return errors.Wrap(err, "Failed to sign TX", logan.F{"tx_being_signed": unsignedOffchainTX})
	}

	err = s.processApprove(ctx, *newRequest, partlySignedOffchainTX, withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to verify Approve Request", logan.F{"partly_signed_offchain_tx_hex": partlySignedOffchainTX})
	}

	return nil
}

func (s *Service) processPreliminaryApprove(ctx context.Context, request horizon.Request, offchainTXHex string) error {
	returnedEnvelope, err := s.sendRequestToVerifier(VerifyPreliminaryApproveURLSuffix, NewApprove(request.ID, request.Hash, offchainTXHex))
	if err != nil {
		return errors.Wrap(err, "Failed to send preliminary Approve to Verify")
	}

	checkErr := checkPreliminaryApproveEnvelope(*returnedEnvelope, request.ID, request.Hash, offchainTXHex)
	if checkErr != "" {
		return errors.Wrap(err, "Envelope returned by Verify is invalid")
	}

	err = s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope")
	}

	s.log.WithField("request", request).WithField("tx_hex", offchainTXHex).
		Debug("Verified PreliminaryApprove successfully.")

	return nil
}

func (s *Service) processApprove(ctx context.Context, request horizon.Request, partlySignedOffchainTX string, withdrawAddress string, withdrawAmount int64) error {
	returnedEnvelope, err := s.sendRequestToVerifier(VerifyApproveURLSuffix, NewApprove(request.ID, request.Hash, partlySignedOffchainTX))
	if err != nil {
		return errors.Wrap(err, "Failed to send Approve to Verify")
	}

	fullySignedOffchainTX, checkErr := s.checkApproveEnvelope(*returnedEnvelope, request.ID, request.Hash, withdrawAddress, withdrawAmount)
	if checkErr != "" {
		return errors.Wrap(err, "Envelope returned by Verify is invalid")
	}

	offchainTXHash, err := s.offchainHelper.SendTX(ctx, fullySignedOffchainTX)
	if err != nil {
		return errors.Wrap(err, "Failed to send fully signed Offchain TX into Offchain network", logan.F{
			"partly_signed_offchain_tx_hex": partlySignedOffchainTX,
			"fully_signed_offchain_tx_hex":  fullySignedOffchainTX,
		})
	}

	err = s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope", logan.F{
			"envelope_returned_by_verify": returnedEnvelope,
			"offchain_sent_tx_hash":       offchainTXHash,
		})
	}

	s.log.WithField("request", request).WithFields(logan.F{
		"sent_offchain_tx_hash": offchainTXHash,
	}).Info("Processed Approve of WithdrawRequest successfully.")

	return nil
}
