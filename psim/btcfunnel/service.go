package btcfunnel

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"time"
)

type BTCClient interface {
	GetWalletBalance() (float64, error)
	SendToAddress(goalAddress string, amount float64) (resultTXHash string, err error)
}

type Service struct {
	config Config
	log    *logan.Entry

	btcClient BTCClient
}

func New(config Config, log *logan.Entry, btcClient BTCClient) *Service {
	return &Service{
		config: config,
		log:    log,

		btcClient: btcClient,
	}
}

// Run will never send errors into the returned channel, however
// when Service stops - it will close the channel.
func (s *Service) Run(ctx context.Context) chan error {
	s.log.Info("Starting.")

	errs := make(chan error)

	go func() {
		app.RunOverIncrementalTimer(ctx, s.log, "btc_funnel_runner", s.funnelBTC, 5 * time.Second, 5 * time.Second)
		close(errs)
	}()

	return errs
}

func (s *Service) funnelBTC(ctx context.Context) error {
	balance, err := s.btcClient.GetWalletBalance()
	if err != nil {
		return errors.Wrap(err, "Failed to get Wallet balance")
	}

	if balance == 0 || balance < s.config.MinFunnelAmount {
		// To less money to funnel.
		return nil
	}

	txHash, err := s.btcClient.SendToAddress(s.config.FunnelAddress, balance)
	if err != nil {
		return errors.Wrap(err, "Failed to send BTC to Address", logan.Field("wallet_confirmed_balance", balance))
	}

	s.log.WithField("funnel_amount", balance).WithField("funnel_tx_hash", txHash).Info("Funneled BTC to the hot wallet.")

	return nil
}
