package withdraw

import (
	"context"

	"time"

	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
)

// ProcessValidPendingRequest knows how to process both TwoStepWithdrawal and Withdraw RequestTypes.
// TODO Make me smaller
func (s *Service) processValidPendingRequest(ctx context.Context, request horizon.Request) error {
	withdrawAddress, err := GetWithdrawalAddress(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get Withdraw Address")
	}
	amount, err := GetWithdrawAmount(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get Withdraw Amount")
	}
	withdrawAmount := s.offchainHelper.ConvertAmount(amount)

	if request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		// TwoStepWithdraw needs PreliminaryApprove first.
		unsignedOffchainTXHex, err := s.offchainHelper.CreateTX(ctx, withdrawAddress, withdrawAmount)
		if err != nil {
			return errors.Wrap(err, "Failed to create Offchain TX")
		}
		if running.IsCancelled(ctx) {
			return nil
		}

		err = s.processPreliminaryApprove(ctx, request, unsignedOffchainTXHex)
		if err != nil {
			return errors.Wrap(err, "Failed to verify Preliminary Approve Request", logan.F{"tx_hex": unsignedOffchainTXHex})
		}
	}

	var newRequest *horizon.Request
	var unsignedOffchainTX string
	if s.verification.Verify {
		newRequest, err = s.requestsConnector.GetRequestByID(request.ID)
		if err != nil {
			return errors.Wrap(err, "Failed to obtain Request from Horizon")
		}

		for newRequest.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
			s.log.WithField("new_request_id", newRequest.ID).
				// TODO sleep period in message - into var or const
				Debugf("WithdrawalRequest still hasn't changed type to Withdraw(%d). Sleeping for 3 seconds.", xdr.ReviewableRequestTypeWithdraw)

			// TODO Incremental
			time.Sleep(3 * time.Second)
			newRequest, err = s.requestsConnector.GetRequestByID(request.ID)
			if err != nil {
				return errors.Wrap(err, "Failed to obtain Request from Horizon")
			}
		}

		unsignedOffchainTX, err = GetTXHex(*newRequest)
		if err != nil {
			return errors.Wrap(err, "Failed to get TX hex from the WithdrawRequest")
		}
	} else {
		// Without verification
		unsignedOffchainTX, err = s.offchainHelper.CreateTX(ctx, withdrawAddress, withdrawAmount)
		if err != nil {
			return errors.Wrap(err, "Failed to create Offchain TX with OffchainHelper")
		}
	}
	if running.IsCancelled(ctx) {
		return nil
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
	var resultEnvelope *xdr.TransactionEnvelope
	var err error

	if s.verification.Verify {
		resultEnvelope, err = s.sendRequestToVerifier(VerifyPreliminaryApproveURLSuffix, NewApprove(request.ID, request.Hash, offchainTXHex))
		if err != nil {
			return errors.Wrap(err, "Failed to send preliminary Approve to Verify")
		}

		checkErr := checkPreliminaryApproveEnvelope(*resultEnvelope, request.ID, request.Hash, offchainTXHex)
		if checkErr != "" {
			return errors.Wrap(err, "Envelope returned by Verify is invalid")
		}
	} else {
		extDetails := ExternalDetails{
			TXHex: offchainTXHex,
		}
		extDetBytes, err := json.Marshal(extDetails)
		if err != nil {
			return errors.Wrap(err, "Failed to marshal ExternalDetails struct into bytes")
		}

		txB64, err := s.xdrbuilder.Transaction(s.verification.SourceKP).Op(xdrbuild.ReviewRequestOp{
			ID:     request.ID,
			Hash:   request.Hash,
			Action: xdr.ReviewRequestOpActionApprove,
			Details: xdrbuild.TwoStepWithdrawalDetails{
				ExternalDetails: string(extDetBytes),
			},
		}).Marshal()
		if err != nil {
			return errors.Wrap(err, "Failed to marshal Approval Transaction")
		}
		fields := logan.F{"raw_tx": txB64}

		resultEnvelope = &xdr.TransactionEnvelope{}
		err = resultEnvelope.Scan(txB64)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal Approval Transaction into Envelope", fields)
		}
	}

	err = s.signAndSubmitEnvelope(ctx, *resultEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope")
	}

	s.log.WithFields(logan.F{
		"request": request,
		"tx_hex":  offchainTXHex,
	}).Debug("Verified PreliminaryApprove successfully.")

	return nil
}

// TODO Refactor me
func (s *Service) processApprove(ctx context.Context, request horizon.Request, partlySignedOffchainTX string, withdrawAddress string, withdrawAmount int64) error {
	var resultEnvelope *xdr.TransactionEnvelope
	var fullySignedOffchainTX string

	if s.verification.Verify {
		envelope, err := s.sendRequestToVerifier(VerifyApproveURLSuffix, NewApprove(request.ID, request.Hash, partlySignedOffchainTX))
		if err != nil {
			return errors.Wrap(err, "Failed to send Approve to Verify")
		}

		var checkErr string
		fullySignedOffchainTX, checkErr = s.checkApproveEnvelope(*envelope, request.ID, request.Hash, withdrawAddress, withdrawAmount)
		if checkErr != "" {
			return errors.Wrap(err, "Envelope returned by Verify is invalid")
		}
		resultEnvelope = envelope
	} else {
		// Working without multisig if no verify
		fullySignedOffchainTX = partlySignedOffchainTX

		txHash, err := s.offchainHelper.GetHash(fullySignedOffchainTX)
		if err != nil {
			return errors.Wrap(err, "Failed to get hash of the Offchain TX")
		}

		extDetails := ExternalDetails{
			TXHex:  fullySignedOffchainTX,
			TXHash: txHash,
		}
		extDetBytes, err := json.Marshal(extDetails)
		if err != nil {
			return errors.Wrap(err, "Failed to marshal ExternalDetails struct into bytes")
		}

		txB64, err := s.xdrbuilder.Transaction(s.verification.SourceKP).Op(xdrbuild.ReviewRequestOp{
			ID:     request.ID,
			Hash:   request.Hash,
			Action: xdr.ReviewRequestOpActionApprove,
			Details: xdrbuild.WithdrawalDetails{
				ExternalDetails: string(extDetBytes),
			},
		}).Marshal()

		if err != nil {
			return errors.Wrap(err, "Failed to marshal Approval Transaction")
		}
		fields := logan.F{"raw_tx": txB64}

		resultEnvelope = &xdr.TransactionEnvelope{}
		err = resultEnvelope.Scan(txB64)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal Approval Transaction into Envelope", fields)
		}
	}

	offchainTXHash, err := s.offchainHelper.SendTX(ctx, fullySignedOffchainTX)
	if err != nil {
		return errors.Wrap(err, "Failed to send fully signed Offchain TX into Offchain network", logan.F{
			"fully_signed_offchain_tx_hex": fullySignedOffchainTX,
		})
	}

	err = s.signAndSubmitEnvelope(ctx, *resultEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope", logan.F{
			"result_ envelope":      resultEnvelope,
			"offchain_sent_tx_hash": offchainTXHash,
		})
	}

	s.log.WithField("request", request).WithFields(logan.F{
		"sent_offchain_tx_hash": offchainTXHash,
	}).Info("Processed Approve of WithdrawRequest successfully.")

	return nil
}
