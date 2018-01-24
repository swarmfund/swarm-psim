package withdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

func (s *Service) processValidPendingRequest(ctx context.Context, request horizon.Request) error {
	withdrawAddress, err := GetWithdrawAddress(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get Withdraw Address")
	}
	withdrawAmount := GetWithdrawAmount(request)

	unsignedOffchainTXHex, err := s.offchainHelper.CreateTX(withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to create Offchain TX")
	}

	err = s.verifyPreliminaryApprove(ctx, request, unsignedOffchainTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to verify Preliminary Approve Request", logan.F{"tx_hex": unsignedOffchainTXHex})
	}

	newRequest, err := ObtainRequest(s.horizon.Client(), request.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain Request from Horizon")
	}

	partlySignedOffchainTXHex, err := s.offchainHelper.SignTX(unsignedOffchainTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to sign TX", logan.F{"tx_being_signed": unsignedOffchainTXHex})
	}

	err = s.verifyApprove(ctx, *newRequest, partlySignedOffchainTXHex, withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to verify Approve Request", logan.F{"partly_signed_offchain_tx_hex": partlySignedOffchainTXHex})
	}

	return nil
}


// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
// FIXME Uncomment
func (s *Service) verifyPreliminaryApprove(ctx context.Context, request horizon.Request, offchainTXHex string) error {
	returnedEnvelope, err := s.sendRequestToVerify(VerifyPreliminaryApproveURLSuffix, NewApprove(request.ID, request.Hash, offchainTXHex))
	if err != nil {
		return errors.Wrap(err, "Failed to send preliminary Approve to Verify")
	}

	checkErr := checkPreliminaryApproveEnvelope(*returnedEnvelope, request.ID, request.Hash, offchainTXHex)
	if checkErr != "" {
		return errors.Wrap(err, "Envelope returned by Verify is invalid")
	}

	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	// FIXME Uncomment
	//err = s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	//if err != nil {
	//	return errors.Wrap(err, "Failed to sign-and-submit Envelope")
	//}

	s.log.WithFields(GetRequestLoganFields("request", request)).WithField("tx_hex", offchainTXHex).
		Debug("Verified PreliminaryApprove successfully.")

	return nil
}

func (s *Service) verifyApprove(ctx context.Context, request horizon.Request, partlySignedOffchainTXHex string, withdrawAddress string, withdrawAmount float64) error {
	returnedEnvelope, err := s.sendRequestToVerify(VerifyApproveURLSuffix, NewApprove(request.ID, request.Hash, partlySignedOffchainTXHex))
	if err != nil {
		return errors.Wrap(err, "Failed to send Approve to Verify")
	}

	fullySignedOffchainTXHex, checkErr := s.checkApproveEnvelope(*returnedEnvelope, request.ID, request.Hash, withdrawAddress, withdrawAmount, s.offchainHelper.GetHotWallerAddress())
	if checkErr != "" {
		return errors.Wrap(err, "Envelope returned by Verify is invalid")
	}

	offchainTXHash, err := s.offchainHelper.SendTX(fullySignedOffchainTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to send fully signed Offchain TX into Offchain network", logan.F{
			"fully_signed_offchain_tx_hex": fullySignedOffchainTXHex,
		})
	}

	err = s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope", logan.F{
			"envelope_returned_by_verify": returnedEnvelope,
			"offchain_sent_tx_hash":       offchainTXHash,
		})
	}

	s.log.WithFields(GetRequestLoganFields("request", request)).WithFields(logan.F{
		"sent_offchain_tx_hash": offchainTXHash,
	}).Info("Verified Approve successfully.")

	return nil
}
