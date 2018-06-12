package listener

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
	Signer             keypair.Full `fig:"signer,required"`
	MixpanelToken      string       `fig:"mixpanel_token"`
	SalesforceUsername string       `fig:"salesforce_username"`
	SalesforcePassword string       `fig:"salesforce_password"`
	TxhistoryCursor    string       `fig:"txhistory_cursor"`
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

const (
	defaultServiceRetryTimeIncrement = 1 * time.Second
	defaultMaxServiceRetryTime       = 30 * time.Second
)

// Run starts dispatching events to analytics services
func (s *Service) Run(ctx context.Context) {
	running.UntilSuccess(ctx, s.logger, conf.ListenerService, s.dispatchEvents, defaultServiceRetryTimeIncrement, defaultMaxServiceRetryTime)
}

func (s *Service) dispatchEvents(ctx context.Context) (bool, error) {
	extractedTxData := s.extractor.Extract(ctx)
	emittedEvents := s.handler.Process(ctx, extractedTxData)
	// s.broadcaster.BroadcastEvents(ctx, emittedEvents)
	for e := range s.broadcaster.BroadcastEvents(ctx, emittedEvents) {
		s.logger.Warn(e)
	}

	return false, nil
}
