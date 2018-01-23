package btcwithdraw

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

func (s *Service) processValidPendingRequest(ctx context.Context, withdrawAddress string, withdrawAmount float64, request horizon.Request) error {
	unsignedBtcTXHex, err := s.createBitcoinTX(withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to create Bitcoin TX")
	}

	err = s.verifyPreliminaryApprove(ctx, request, unsignedBtcTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to verify Preliminary Approve Request", logan.F{"btc_tx_hex": unsignedBtcTXHex})
	}

	newRequest, err := withdraw.ObtainRequest(s.horizon.Client(), request.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain Request from Horizon")
	}

	partlySignedBtcTXHex, err := s.btcClient.SignAllTXInputs(unsignedBtcTXHex, s.config.HotWalletScriptPubKey, s.config.HotWalletRedeemScript, s.config.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to sign BTC TX", logan.F{"btc_tx_being_signed": unsignedBtcTXHex})
	}

	err = s.verifyApprove(ctx, *newRequest, partlySignedBtcTXHex, withdrawAddress, withdrawAmount, s.btcClient.GetNetParams())
	if err != nil {
		return errors.Wrap(err, "Failed to verify Approve Request", logan.F{"partly_signed_btc_tx_hex": partlySignedBtcTXHex})
	}

	return nil
}

func (s *Service) createBitcoinTX(withdrawAddress string, withdrawAmount float64) (string, error) {
	txHex, err := s.btcClient.CreateAndFundRawTX(withdrawAddress, withdrawAmount, s.config.HotWalletAddress)
	if err != nil {
		if errors.Cause(err) == bitcoin.ErrInsufficientFunds {
			return "", errors.Wrap(err, "Could not create raw TX - not enough BTC on hot wallet")
		}

		return "", errors.Wrap(err, "Failed to create raw TX")
	}

	return txHex, nil
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
func (s *Service) verifyPreliminaryApprove(ctx context.Context, request horizon.Request, btcTXHex string) error {
	returnedEnvelope, err := s.sendRequestToVerify(withdraw.VerifyPreliminaryApproveURLSuffix, withdraw.NewApprove(request.ID, request.Hash, btcTXHex))
	if err != nil {
		return errors.Wrap(err, "Failed to send preliminary Approve to Verify")
	}

	checkErr := checkPreliminaryApproveEnvelope(*returnedEnvelope, request.ID, request.Hash, btcTXHex)
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

	s.log.WithFields(withdraw.GetRequestLoganFields("request", request)).WithField("tx_hex", btcTXHex).
		Info("Sent Approve to Verify successfully.")

	return nil
}

func (s *Service) verifyApprove(ctx context.Context, request horizon.Request, partlySignedBtcTXHex string, withdrawAddress string, withdrawAmount float64, netParams *chaincfg.Params) error {
	returnedEnvelope, err := s.sendRequestToVerify(withdraw.VerifyApproveURLSuffix, withdraw.NewApprove(request.ID, request.Hash, partlySignedBtcTXHex))
	if err != nil {
		return errors.Wrap(err, "Failed to send Approve to Verify")
	}

	fullySignedBtcTXHex, checkErr := checkApproveEnvelope(*returnedEnvelope, request.ID, request.Hash, withdrawAddress, withdrawAmount, s.config.HotWalletAddress, netParams)
	if checkErr != "" {
		return errors.Wrap(err, "Envelope returned by Verify is invalid")
	}

	btcTXHash, err := s.btcClient.SendRawTX(fullySignedBtcTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to send fully signed BTC TX into Bitcoin network")
	}

	err = s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope")
	}

	s.log.WithFields(withdraw.GetRequestLoganFields("request", request)).WithFields(logan.F{
		"fully_signed_tx_hex": fullySignedBtcTXHex,
		"btc_tx_hash":         btcTXHash,
	}).Info("Sent Approve to Verify successfully.")

	return nil
}
