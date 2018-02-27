package verifier

import (
	"context"
	"net"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/tokend/keypair"
)

type Verifier interface {
	GetOperationType() xdr.OperationType
	VerifyEnvelope(xdr.TransactionEnvelope) (verifyErr, err error)
}

type DiscoveryClient interface {
	EnsureServiceRegistered(service *discovery.Service) (bool, error)
	Service(registration *discovery.ServiceRegistration) *discovery.Service
}

type Service struct {
	serviceName string
	log         *logan.Entry

	verifier   Verifier
	xdrbuilder *xdrbuild.Builder
	signer     keypair.Full
	listener   net.Listener

	discoveryRegisterPeriod time.Duration
	discovery               DiscoveryClient
	discoveryService        *discovery.Service
}

func New(
	serviceName string,
	discoveryID string,
	log *logan.Entry,
	verifier Verifier,
	builder *xdrbuild.Builder,
	signer keypair.Full,
	listener net.Listener,
	discoveryClient DiscoveryClient) *Service {

	discoveryRegisterPeriod := 5 * time.Second

	return &Service{
		serviceName: serviceName,
		log:         log.WithField("service", serviceName),

		verifier:   verifier,
		xdrbuilder: builder,
		signer:     signer,
		listener:   listener,

		discoveryRegisterPeriod: discoveryRegisterPeriod,
		discovery:               discoveryClient,
		discoveryService: discoveryClient.Service(&discovery.ServiceRegistration{
			Name:            serviceName,
			ID:              discoveryID,
			TTL:             2 * discoveryRegisterPeriod,
			DeregisterAfter: 3 * discoveryRegisterPeriod,

			Host: "http://" + listener.Addr().String(),
		}),
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
