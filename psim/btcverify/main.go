package btcverify

import (
	"context"
	"fmt"
	"net"
	"time"

	"sync"

	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	setupFn := func(ctx context.Context) (utils.Service, error) {
		serviceConfig := Config{
			Host:        "localhost",
			ServiceName: conf.ServiceBTCVerify,
		}

		globalConfig := app.Config(ctx)
		err := figure.
			Out(&serviceConfig).
			From(globalConfig.Get(conf.ServiceBTCVerify)).
			With(figure.BaseHooks, utils.CommonHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceBTCVerify))
		}

		log := ctx.Value(app.CtxLog).(*logan.Entry).WithField("service", conf.ServiceBTCVerify)

		discoveryClient, err := globalConfig.Discovery()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get discovery client")
		}

		horizonConnector, err := globalConfig.Horizon()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get Horizon connector")
		}

		listener, err := ape.Listener(serviceConfig.Host, serviceConfig.Port)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init listener")
		}

		btcClient, err := globalConfig.Bitcoin()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get Bitcoin client")
		}

		return newService(ctx, serviceConfig, log, discoveryClient, listener, horizonConnector, btcClient), nil
	}

	app.RegisterService(conf.ServiceBTCVerify, setupFn)
}

// BTCClient must be implemented by a BTC Client to pass into Service for creation.
type btcClient interface {
	IsTestnet() bool
	GetBlockByHash(blockHash string) (*btc.Block, error)
}

type Service struct {
	ServiceID string

	ctx      context.Context
	config   Config
	log      *logan.Entry
	errors   chan error
	listener net.Listener
	horizon  *horizon.Connector

	discovery        *discovery.Client
	discoveryService *discovery.Service

	btcClient btcClient
}

func newService(ctx context.Context, config Config, log *logan.Entry, discovery *discovery.Client, listener net.Listener,
	horizon *horizon.Connector, btcClient btcClient) *Service {

	return &Service{
		ServiceID: utils.GenerateToken(),
		ctx:       ctx,
		config:    config,
		log:       log,
		errors:    make(chan error),
		listener:  listener,
		horizon:   horizon,

		discovery: discovery,

		btcClient: btcClient,
	}
}

//Run starts all runners in separate goroutines and creates routine, which waits for all of the runners to return.
//Once all runners returned - Errors channel will be closed.
//Implements utils.Service.
func (s *Service) Run() chan error {
	runners := []func(){
		s.registerInDiscovery,
		s.serveAPI,
	}

	go func() {
		wg := sync.WaitGroup{}

		for _, runner := range runners {
			ohigo := runner
			wg.Add(1)

			go func() {
				ohigo()
				wg.Done()
			}()
		}

		wg.Wait()
		close(s.errors)
	}()

	return s.errors
}

func (s *Service) registerInDiscovery() {
	s.discoveryService = s.discovery.Service(&discovery.ServiceRegistration{
		Name: s.config.ServiceName,
		ID:   s.ServiceID,
		Host: fmt.Sprintf("http://%s", s.listener.Addr().String()),
	})

	// FIXME Select from ticker and ctx.Done() simultaneously
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		if app.IsCanceled(s.ctx) {
			return
		}

		err := s.discovery.RegisterServiceSync(s.discoveryService)
		if err != nil {
			s.errors <- errors.Wrap(err, "discovery error")
			continue
		}
	}
}

// VerifyRequest is struct, which is used to parse requests coming from btcsupervisor,
// btcsupervisor should use this type for requests to this verifier.
type VerifyRequest struct {
	Envelope string `json:"envelope"`

	BlockHash string `json:"block_hash"`
	TXHash    string `json:"tx_hash"`
	OutIndex  int    `json:"out_index"`
}
