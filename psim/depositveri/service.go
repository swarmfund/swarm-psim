package depositveri

import (
	"context"
	"net"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/deposit"
	"gitlab.com/tokend/keypair"
)

type Service struct {
	serviceName string

	log *logan.Entry

	signer             keypair.Full
	lastBlocksNotWatch uint64

	// TODO Interface
	horizon    *horizon.Connector
	xdrbuilder *xdrbuild.Builder
	listener   net.Listener

	discoveryRegisterPeriod time.Duration
	// TODO Interface
	discovery        *discovery.Client
	discoveryService *discovery.Service

	offchainHelper deposit.OffchainHelper
}

func New(
	serviceName string,
	log *logan.Entry,
	signer keypair.Full,
	lastBlocksNotWatch uint64,
	horizon *horizon.Connector,
	builder *xdrbuild.Builder,
	listener net.Listener,
	discoveryClient *discovery.Client,
	offchainHelper deposit.OffchainHelper) *Service {

	discoveryRegisterPeriod := 5 * time.Second

	return &Service{
		serviceName: serviceName,

		log: log.WithField("service", serviceName),

		signer:             signer,
		lastBlocksNotWatch: lastBlocksNotWatch,

		horizon:    horizon,
		xdrbuilder: builder,
		listener:   listener,

		discoveryRegisterPeriod: discoveryRegisterPeriod,
		discovery:               discoveryClient,
		discoveryService: discoveryClient.Service(&discovery.ServiceRegistration{
			Name:            serviceName,
			ID:              "my_awesome_super_duper_random_id_deposit",
			TTL:             2 * discoveryRegisterPeriod,
			DeregisterAfter: 3 * discoveryRegisterPeriod,

			Host: "http://" + listener.Addr().String(),
		}),

		offchainHelper: offchainHelper,
	}
}

// TODO Comment
func (s *Service) Run(ctx context.Context) {
	// TODO Wait for acquireLeadershipEndlessly on shutdown
	go app.RunOverIncrementalTimer(ctx, s.log, s.serviceName+"_discovery_reregisterer", s.ensureServiceInDiscoveryOnce,
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
