package eventsubmitter

import (
	"context"
	"time"

	"gitlab.com/swarmfund/psim/psim/conf"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/keypair"
)

// ServiceConfig holds signer for horizon connector and some data for targets
type ServiceConfig struct {
	Signer          keypair.Full `fig:"signer,required"`
	TxHistoryCursor string       `fig:"txhistory_cursor"`
	Targets         []string     `fig:"targets,required"`
}

// Service consists config, logger, broadcaster and dependent components - extractor and handler
type Service struct {
	config      ServiceConfig
	extractor   Extractor
	handler     Handler
	broadcaster Broadcaster
	logger      *logan.Entry
}

// NewService constructs a Service from provided fields
func NewService(config ServiceConfig, extractor Extractor, handler Handler, broadcaster Broadcaster, log *logan.Entry) *Service {
	return &Service{
		config:      config,
		extractor:   extractor,
		handler:     handler,
		broadcaster: broadcaster,
		logger:      log,
	}
}

// Run starts dispatching events to analytics services
func (s *Service) Run(ctx context.Context) {
	s.logger.Info("starting")
	running.UntilSuccess(ctx, s.logger, conf.EventSubmitterService, s.dispatchEvents, 1*time.Second, 30*time.Second)
}

func (s *Service) dispatchEvents(ctx context.Context) (bool, error) {
	extractedTxData := s.extractor.Extract(ctx)
	emittedEvents := s.handler.Process(ctx, extractedTxData)
	s.broadcaster.BroadcastEvents(ctx, emittedEvents)
	return true, nil
}
