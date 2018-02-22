package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/issuance"
)

func (s *Service) processIssuance(ctx context.Context, accountAddress string, issuanceOpt issuance.RequestOpt) error {
	tx := issuance.CraftIssuanceTX(issuanceOpt, s.builder, s.source, s.signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal TX into Envelope")
	}

	logger := s.log.WithFields(logan.F{
		"account_address": accountAddress,
		"issuance":        issuanceOpt,
	})

	ok, err := issuance.SubmitEnvelope(ctx, envelope, s.txSubmitter)
	if err != nil {
		return errors.Wrap(err, "Failed to submit IssuanceRequest TX Envelope to Horizon")
	}

	if ok {
		logger.Info("CoinEmissionRequest was sent successfully.")
	} else {
		logger.Debug("Reference duplication - already processed Deposit, skipping.")
	}
	return nil
}
