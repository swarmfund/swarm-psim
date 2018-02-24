package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/issuance"
)

func (s *Service) submitIssuance(ctx context.Context, accountAddress, balanceID string) (bool, error) {
	issuanceOpt := issuance.RequestOpt{
		Reference: s.buildReference(accountAddress),
		Receiver:  balanceID,
		Asset:     s.config.Asset,
		Amount:    s.config.Amount,
		Details:   `{"cause": "airdrop"}`,
	}

	tx := issuance.CraftIssuanceTX(issuanceOpt, s.builder, s.config.Source, s.config.Signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return false, errors.Wrap(err, "Failed to marshal TX into Envelope")
	}

	logger := s.log.WithFields(logan.F{
		"account_address": accountAddress,
		"issuance":        issuanceOpt,
	})

	ok, err := issuance.SubmitEnvelope(ctx, envelope, s.txSubmitter)
	if err != nil {
		return false, errors.Wrap(err, "Failed to submit IssuanceRequest TX Envelope to Horizon")
	}

	if ok {
		logger.Info("CoinEmissionRequest was sent successfully.")
		return true, nil
	} else {
		logger.Info("Reference duplication - already processed Deposit, skipping.")
		return false, nil
	}
}

func (s *Service) buildReference(accountAddress string) string {
	result := accountAddress + "-airdrop" // accountAddress should be 56 runes length

	// Just in case.
	if len(result) > 64 {
		result = result[:64]
	}

	// Just in case.
	if len(result) < 64 {
		filler := "----------------------------------------------------------------" // len = 64
		result = result + filler[:64-len(result)]
	}

	return result
}
