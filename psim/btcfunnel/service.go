package btcfunnel

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"time"
)

// BTCClient is the interface to be implemented by a
// Bitcoin client to parametrize the Service.
type BTCClient interface {
	GetWalletBalance() (float64, error)
	SendToAddress(goalAddress string, amount float64) (resultTXHash string, err error)
}

// Service implements utils.Service to be registered in the app.
type Service struct {
	config Config
	log    *logan.Entry

	btcClient BTCClient
}

// New is constructor for btcfunnel Service.
func New(config Config, log *logan.Entry, btcClient BTCClient) *Service {
	return &Service{
		config: config,
		log:    log,

		btcClient: btcClient,
	}
}

// Run will return closed channel and only when work is finished.
func (s *Service) Run(ctx context.Context) chan error {
	s.log.Info("Starting.")

	app.RunOverIncrementalTimer(ctx, s.log, "btc_funnel_runner", s.funnelBTC, 5 * time.Second, 5 * time.Second)

	errs := make(chan error)
	close(errs)
	return errs
}

func (s *Service) funnelBTC(ctx context.Context) error {
	balance, err := s.btcClient.GetWalletBalance()
	if err != nil {
		return errors.Wrap(err, "Failed to get Wallet balance")
	}

	if balance == 0 || balance < s.config.MinFunnelAmount {
		// Too little money to funnel.
		return nil
	}

	txHash, err := s.btcClient.SendToAddress(s.config.FunnelAddress, balance)
	if err != nil {
		return errors.Wrap(err, "Failed to send BTC to Address", logan.Field("wallet_confirmed_balance", balance))
	}

	s.log.WithField("funnel_amount", balance).WithField("funnel_tx_hash", txHash).Info("Funneled BTC to the hot wallet.")

	return nil
}
