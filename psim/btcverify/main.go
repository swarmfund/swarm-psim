package btcverify

import (
	"context"
	"fmt"
	"net"
	"time"

	"sync"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func init() {
	setupFn := func(ctx context.Context) (app.Service, error) {
		serviceConfig := Config{
			Host:        "localhost",
			ServiceName: conf.ServiceBTCVerify,
		}

		globalConfig := app.Config(ctx)
		err := figure.
			Out(&serviceConfig).
			From(globalConfig.GetRequired(conf.ServiceBTCVerify)).
			With(figure.BaseHooks, utils.CommonHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceBTCVerify))
		}

		log := app.Log(ctx).WithField("service", conf.ServiceBTCVerify)

		listener, err := ape.Listener(serviceConfig.Host, serviceConfig.Port)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init listener")
		}

		return newService(serviceConfig, log, globalConfig.Discovery(), listener, globalConfig.Horizon(), globalConfig.Bitcoin()), nil
	}

	app.RegisterService(conf.ServiceBTCVerify, setupFn)
}

// BTCClient must be implemented by a BTC Client to pass into Service for creation.
type btcClient interface {
	GetBlockByHash(blockHash string) (*btcutil.Block, error)
}

type Service struct {
	ServiceID string

	config   Config
	log      *logan.Entry
	listener net.Listener
	horizon  *horizon.Connector

	discovery        *discovery.Client
	discoveryService *discovery.Service

	btcClient btcClient
}

func newService(config Config, log *logan.Entry, discovery *discovery.Client, listener net.Listener,
	horizon *horizon.Connector, btcClient btcClient) *Service {

	return &Service{
		ServiceID: utils.GenerateToken(),
		config:    config,
		log:       log,
		listener:  listener,
		horizon:   horizon,

		discovery: discovery,

		btcClient: btcClient,
	}
}

// Run starts all runners in separate goroutines and creates routine, which waits for all of the runners to return.
// Once all runners returned - this method will finish.
// Implements app.Service.
func (s *Service) Run(ctx context.Context) {
	runners := []func(context.Context){
		s.registerInDiscovery,
		s.serveAPI,
	}

	wg := sync.WaitGroup{}

	for _, runner := range runners {
		ohigo := runner
		wg.Add(1)

		go func() {
			ohigo(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (s *Service) registerInDiscovery(ctx context.Context) {
	s.discoveryService = s.discovery.Service(&discovery.ServiceRegistration{
		Name: s.config.ServiceName,
		ID:   s.ServiceID,
		Host: fmt.Sprintf("http://%s", s.listener.Addr().String()),
	})

	// FIXME Select from ticker and ctx.Done() simultaneously
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		if app.IsCanceled(ctx) {
			return
		}

		err := s.discovery.RegisterServiceSync(s.discoveryService)
		if err != nil {
			s.log.WithError(err).Error("discovery error")
			continue
		}
	}
}
