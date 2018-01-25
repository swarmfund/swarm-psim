package bearer

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

// Service is a main structure for bearer runner,
// implements `app.Service` interface.
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
		logger:  log.WithField("service", conf.ServiceBearer),
	}
}

// Run will returns only when work is finished.
func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting.")

	app.RunOverIncrementalTimer(
		ctx,
		s.logger,
		conf.ServiceBearer,
		s.sendOperations,
		0,
		s.config.AbnormalPeriod)
}

// sendOperations is create and submit operations.
func (s *Service) sendOperations(ctx context.Context) error {
	err := s.checkSaleState()
	if err == nil {
		s.logger.Info("Operation submitted")
		return nil
	}

	if err != errorNoSales {
		return errors.Wrap(err, "can not to submit checkSaleState operation")
	}

	tm := time.NewTimer(s.config.SleepPeriod)
	select {
	case <-ctx.Done():
		return nil
	case <-tm.C:
		return nil
	}
}
