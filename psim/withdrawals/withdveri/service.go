package withdveri

import (
	"context"
	"net"
	"sync"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/swarmfund/psim/psim/withdrawals/withdraw"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/regources"
)

type RequestsConnector interface {
	GetRequestByID(requestID uint64) (*regources.ReviewableRequest, error)
}

type Service struct {
	serviceName string

	log *logan.Entry

	sourceKP keypair.Address
	signerKP keypair.Full

	requestsConnector RequestsConnector
	xdrbuilder        *xdrbuild.Builder
	listener          net.Listener

	discoveryRegisterPeriod time.Duration
	discovery               *discovery.Client
	discoveryService        *discovery.Service

	offchainHelper withdraw.CommonOffchainHelper
}

func New(
	serviceName string,
	log *logan.Entry,
	sourceKP keypair.Address,
	signerKP keypair.Full,
	requestsConnector RequestsConnector,
	builder *xdrbuild.Builder,
	listener net.Listener,
	discoveryClient *discovery.Client,
	offchainHelper withdraw.CommonOffchainHelper) *Service {

	discoveryRegisterPeriod := 5 * time.Second

	return &Service{
		serviceName: serviceName,

		log: log.WithField("service", serviceName),

		sourceKP: sourceKP,
		signerKP: signerKP,

		requestsConnector: requestsConnector,
		xdrbuilder:        builder,
		listener:          listener,

		discoveryRegisterPeriod: discoveryRegisterPeriod,
		discovery:               discoveryClient,
		discoveryService: discoveryClient.Service(&discovery.ServiceRegistration{
			Name:            serviceName,
			ID:              utils.GenerateToken(),
			TTL:             2 * discoveryRegisterPeriod,
			DeregisterAfter: 3 * discoveryRegisterPeriod,

			Host: "http://" + listener.Addr().String(),
		}),

		offchainHelper: offchainHelper,
	}
}

// TODO Comment
func (s *Service) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		running.WithBackOff(ctx, s.log, "discovery_reregisterer", s.ensureServiceInDiscoveryOnce,
			s.discoveryRegisterPeriod, s.discoveryRegisterPeriod/2, s.discoveryRegisterPeriod*10)
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
