package btcfunnel

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"fmt"
)

func (s *Service) funnelBTC(ctx context.Context) error {
	balance, err := s.btcClient.GetWalletBalance(false)
	if err != nil {
		return errors.Wrap(err, "Failed to get Wallet balance")
	}

	if balance == 0 || balance < s.config.MinFunnelAmount {
		// Too little money to funnel.
		return nil
	}

	fields := logan.F{
		"funnel_amount": balance,
	}

	balanceWithWatchOnly, err := s.btcClient.GetWalletBalance(true)
	if err != nil {
		return errors.Wrap(err, "Failed to get Wallet balance with watch only", fields)
	}
	hotBalance := balanceWithWatchOnly - balance
	fields["hot_balance"] = fmt.Sprintf("%.8f", hotBalance)

	// Prepare amounts
	amountToHot, amountToCold := s.countAmounts(balance, hotBalance)
	fields["amount_to_hot"] = fmt.Sprintf("%.8f", amountToHot)
	fields["amount_to_cold"] = fmt.Sprintf("%.8f", amountToCold)

	addrToAmount := make(map[string]float64)
	if amountToHot > 0 {
		addrToAmount[s.config.HotAddress] = amountToHot
	}
	if amountToCold > 0 {
		addrToAmount[s.config.ColdAddress] = amountToCold
	}

	txHash, err := s.btcClient.SendMany(addrToAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to send BTC to Hot/Cold Addresses", fields)
	}
	fields["funnel_tx_hash"] = txHash

	s.log.WithFields(fields).Info("Funneled BTC to the Hot/Cold Address(es).")

	return nil
}

// TODO Cover with tests.
func (s *Service) countAmounts(sendingAmount, hotBalance float64) (amountToHot, amountToCold float64) {
	availableToHot := s.config.MaxHotStock - hotBalance

	if availableToHot < s.config.DustOutputLimit {
		// Already enough BTC on the Hot Address - sending everything to the Cold.
		return 0, sendingAmount
	}

	if sendingAmount - availableToHot > s.config.DustOutputLimit {
		// After fulfilling of the Hot, still have more than Dust to be sent to the Cold.
		return availableToHot, sendingAmount - availableToHot
	}

	// Whole money to the Hot.
	return sendingAmount, 0
}
