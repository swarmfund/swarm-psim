package reporter

import (
	"context"
	horizon "horizon-connector"
	"time"

	"gitlab.com/swarmfund/psim/psim/conf"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/keypair"
)

type ServiceConfig struct {
	Signer keypair.Full `fig:"signer,required"`
	// TODO tx threshold
	// TODO asset_code
}

type Service struct {
	config      ServiceConfig
	horizon     *horizon.Connector
	broadcaster Broadcaster
	logger      *logan.Entry
}

func NewService(config ServiceConfig, horizon *horizon.Connector, broadcaster Broadcaster, log *logan.Entry) *Service {
	return &Service{
		horizon:     horizon,
		broadcaster: broadcaster,
		logger:      log,
	}
}

const (
	defaultServiceRetryTimeIncrement = 1 * time.Second
	defaultMaxServiceRetryTime       = 30 * time.Second
)

func (s *Service) Run(ctx context.Context) {
	running.UntilSuccess(ctx, s.logger, conf.ListenerService, s.dispatchEvents, defaultServiceRetryTimeIncrement, defaultMaxServiceRetryTime)
}

type ProcessedItem string

func (s *Service) dispatchEvents(ctx context.Context) (bool, error) {
	emittedEvents := make(chan ProcessedItem)
	emittedEvents <- s.horizon.System().Balances()
	close(emittedEvents)
	s.broadcaster.BroadcastEvents(ctx, emittedEvents)
	return true, nil
}
