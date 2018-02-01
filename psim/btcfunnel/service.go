package btcfunnel

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
	"time"
)

// BTCClient is the interface to be implemented by a
// Bitcoin client to parametrize the Service.
type BTCClient interface {
	GetWalletBalance(includeWatchOnly bool) (float64, error)
	SendMany(addrToAmount map[string]float64) (resultTXHash string, err error)
}

// Service implements app.Service to be registered in the app.
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

// Run is implementation of app.Service, Run is called by the app.
// Run will return only when work is finished.
func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")
	app.RunOverIncrementalTimer(ctx, s.log, "btc_funnel_runner", s.funnelEverythingFromSmallAddresses, 5 * time.Second, 5 * time.Second)
}
