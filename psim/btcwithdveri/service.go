package btcwithdveri

import (
	"context"
	"net"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

type BTCClient interface {
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript string, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
	GetNetParams() *chaincfg.Params
}

type Service struct {
	log    *logan.Entry
	config Config

	horizon    *horizon.Connector
	xdrbuilder *xdrbuild.Builder
	listener   net.Listener
	btcClient  BTCClient

	discoveryRegisterPeriod time.Duration
	discovery               *discovery.Client
	discoveryService        *discovery.Service
}

func New(log *logan.Entry, config Config,
	horizon *horizon.Connector, builder *xdrbuild.Builder, btc BTCClient, listener net.Listener, discoveryClient *discovery.Client) *Service {

	discoveryRegisterPeriod := 5 * time.Second

	return &Service{
		log:    log.WithField("service", conf.ServiceBTCWithdrawVerify),
		config: config,

		horizon:    horizon,
		xdrbuilder: builder,
		listener:   listener,
		btcClient:  btc,

		discoveryRegisterPeriod: discoveryRegisterPeriod,
		discovery:               discoveryClient,
		discoveryService: discoveryClient.Service(&discovery.ServiceRegistration{
			Name:            conf.ServiceBTCWithdrawVerify,
			ID:              "my_awesome_super_duper_random_id",
			TTL:             2 * discoveryRegisterPeriod,
			DeregisterAfter: 3 * discoveryRegisterPeriod,

			Host: "http://" + listener.Addr().String(),
		}),
	}
}

// TODO Comment
func (s *Service) Run(ctx context.Context) {
	// TODO Wait for acquireLeadershipEndlessly on shutdown
	go app.RunOverIncrementalTimer(ctx, s.log, "btc_withdraw_verify_discovery_reregisterer", s.ensureServiceInDiscoveryOnce,
		s.discoveryRegisterPeriod, s.discoveryRegisterPeriod/2)

	s.serveAPI(ctx)
}

func (s *Service) ensureServiceInDiscoveryOnce(ctx context.Context) error {
	_, err := s.discovery.EnsureServiceRegistered(s.discoveryService)
	if err != nil {
		return errors.Wrap(err, "Failed to ensure service registered in Discovery")
	}
	return nil
}
