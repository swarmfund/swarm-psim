package notifier

import (
	"context"

	"sync"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/sse-go"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

type Service struct {
	*Config
	horizon *horizon.Connector
	sender  *notificator.Connector
	logger  *logan.Entry
	sse     *sse.Listener
}

// New returns new instance of the Service service.
func New(
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
	}
}

// Run start service executing
// returns chan of processing errors.
func (s *Service) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	// TODO Make runners return error
	// TODO Consider running runners over Incremental timer
	enabledRunners := []func(ctx context.Context){
		s.checkAssetsIssuanceAmount,
		s.listenOperations,
		s.servePProfAPI,
	}

	for _, fn := range enabledRunners {
		serviceRunner := fn
		wg.Add(1)
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					s.logger.WithError(errors.FromPanic(rec)).Error("runner panicked")
				}
				wg.Done()
			}()
			serviceRunner(ctx)
		}()
	}

	wg.Wait()
}
