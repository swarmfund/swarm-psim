package notifier

import (
	"context"

	"sync"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/sse-go"
	"gitlab.com/swarmfund/horizon-connector"
)

type Service struct {
	*Config
	horizon *horizon.Connector
	sender  *notificator.Connector
	logger  *logan.Entry
	errors  chan error
	sse     *sse.Listener

	// teardown
	ctx    context.Context
	cancel context.CancelFunc
}

// New returns new instance of the Service service.
func New(
	ctx context.Context,
	config *Config,
	sender *notificator.Connector,
	horizonConn *horizon.Connector,
	logger *logan.Entry,
) *Service {
	return &Service{
		Config:  config,
		horizon: horizonConn,
		logger:  logger,
		sender:  sender,
		ctx:     ctx,
	}
}

// Run start service executing
// returns chan of processing errors.
func (s *Service) Run() chan error {
	s.errors = make(chan error)

	wg := sync.WaitGroup{}
	enabledServices := []func(ctx context.Context){
		s.checkAssetsIssuanceAmount,
		s.listenOperations,
		s.servePProfAPI,
	}

	for _, fn := range enabledServices {
		serviceRunner := fn
		wg.Add(1)
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					s.errors <- errors.FromPanic(rec)
				}
				wg.Done()
			}()
			serviceRunner(s.ctx)
		}()
	}

	go func() {
		defer func() {
			close(s.errors)
		}()
		wg.Wait()
	}()

	return s.errors
}
