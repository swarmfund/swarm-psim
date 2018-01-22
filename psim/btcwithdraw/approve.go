package btcwithdraw

import (
	"context"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

func (s *Service) processValidPendingRequest(ctx context.Context, withdrawAddress string, withdrawAmount float64, request horizon.Request) error {
	btcTXHex, err := s.createBitcoinTX(withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to create Bitcoin TX")
	}

	err = s.verifyApprove(ctx, request, btcTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to verify Approve Request", logan.F{"btc_tx_hex": btcTXHex})
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

func (s *Service) verifyApprove(ctx context.Context, request horizon.Request, btcTXHex string) error {
	returnedEnvelope, err := s.sendRequestToVerify(withdraw.NewApprove(request.ID, request.Hash, btcTXHex))
	if err != nil {
		return errors.Wrap(err, "Failed to send preliminary Approve to Verify")
	}

	checkErr := checkApproveEnvelope(*returnedEnvelope, request.ID, request.Hash, btcTXHex)
	if checkErr != "" {
		return errors.Wrap(err, "Envelope returned by Verify is invalid")
	}

	err = s.signAndSubmitEnvelope(ctx, *returnedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to sign-and-submit Envelope")
	}

	s.log.WithFields(withdraw.GetRequestLoganFields("request", request)).WithField("tx_hex", btcTXHex).
		Info("Sent Approve to Verify successfully.")

	return nil
}
