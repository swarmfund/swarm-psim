package btcwithdveri

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/conf"
	"net"
	"gitlab.com/distributed_lab/discovery-go"
	"time"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type BTCClient interface {
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript *string, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
	IsTestnet() bool
}

type Service struct {
	log    *logan.Entry
	config Config

	horizon   *horizon.Connector
	btcClient BTCClient
	listener  net.Listener

	discoveryRegisterPeriod time.Duration
	discovery *discovery.Client
	discoveryService *discovery.Service
}

func New(log *logan.Entry, config Config,
	horizon *horizon.Connector, btc BTCClient, listener net.Listener, discoveryClient *discovery.Client) *Service {

	discoveryRegisterPeriod := 5 * time.Second

	return &Service{
		log:    log.WithField("service", conf.ServiceBTCWithdrawVerify),
		config: config,

		horizon:   horizon,
		btcClient: btc,
		listener:  listener,

		discoveryRegisterPeriod: discoveryRegisterPeriod,
		discovery: discoveryClient,
		discoveryService : discoveryClient.Service(&discovery.ServiceRegistration{
			Name: conf.ServiceBTCWithdrawVerify,
			ID:   "my_awesome_super_duper_random_id",
			TTL:  2 * discoveryRegisterPeriod,
			DeregisterAfter: 3 * discoveryRegisterPeriod,

			Host: "http://" + listener.Addr().String(),
		}),
	}
}

func (s *Service) Run(ctx context.Context) chan error {
	// TODO Wait for acquireLeadershipEndlessly on shutdown
	go app.RunOverIncrementalTimer(ctx, s.log, "btc_withdraw_verify_discovery_reregisterer", s.ensureServiceInDiscoveryOnce,
		s.discoveryRegisterPeriod, s.discoveryRegisterPeriod / 2)

	s.serveAPI(ctx)

	errs := make(chan error)
	close(errs)
	return errs
}

func (s *Service) ensureServiceInDiscoveryOnce(ctx context.Context) error {
	_, err := s.discovery.EnsureServiceRegistered(s.discoveryService)
	if err != nil {
		errors.Wrap(err, "Failed to ensure service registered in Discovery")
	}
	return nil
}
