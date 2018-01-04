package bearer

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/app"
)

type Service struct {
	config  Config
	horizon *horizon.Connector
	logger  *logan.Entry
	errors  chan error
}

// New is constructor for bearer Service.
func New(config Config, log *logan.Entry, connector *horizon.Connector) *Service {
	return &Service{
		config:  config,
		horizon: connector,
		logger:  log,
	}
}

// Run will return closed channel and only when work is finished.
func (s *Service) Run(ctx context.Context) chan error {
	s.logger.Info("Starting.")

	app.RunOverIncrementalTimer(ctx, s.logger, "operation_bearer", s.operationBearer, s.config.Period, 2*s.config.Period)

	errs := make(chan error)
	close(errs)
	return errs
}

func (s *Service) operationBearer(ctx context.Context) error {
	return nil
}
