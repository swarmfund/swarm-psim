package verifier

import (
	"context"
	"net"
	"time"

	"sync"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/keypair"
	"gitlab.com/distributed_lab/running"
)

type Verifier interface {
	Run(ctx context.Context)
	// GetOperationType must return the OperationType Verifier expects in the TX.
	GetOperationType() xdr.OperationType
	// On calling of this method it's guaranteed that provided Envelope has exactly 1 Operation
	// and that the Type of Operation equals to return of the GetOperationType interface method.
	VerifyOperation(xdr.TransactionEnvelope) (verifyErr, err error)
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

	discoveryID             string
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

		discoveryID:             discoveryID,
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
	s.log.WithFields(logan.F{
		"discovery_id":              s.discoveryID,
		"discovery_register_period": s.discoveryRegisterPeriod,
	}).Info("Starting general verifying service with provided Verifier.")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.verifier.Run(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		running.WithBackOff(ctx, s.log, "discovery_reregisterer", s.ensureServiceInDiscoveryOnce,
			s.discoveryRegisterPeriod, s.discoveryRegisterPeriod/2, time.Minute)
		wg.Done()
	}()

	s.serveAPI(ctx)
	wg.Wait()
}

func (s *Service) ensureServiceInDiscoveryOnce(ctx context.Context) error {
	_, err := s.discovery.EnsureServiceRegistered(s.discoveryService)
	if err != nil {
		return errors.Wrap(err, "Failed to ensure service registered in Discovery")
	}
	return nil
}
