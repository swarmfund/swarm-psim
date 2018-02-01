package bearer

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

type SaleStateCheckerInterface interface {
	// Return sales from core DB
	GetSales() ([]horizon.Sale, error)
	GetHorizonInfo() (info *horizon.Info, err error)
	BuildTx(info *horizon.Info, saleID uint64) (string, error)
	SubmitTx(ctx context.Context, envelope string) horizon.SubmitResult
}

// Service is a main structure for bearer runner,
// implements `app.Service` interface.
type Service struct {
	config  Config
	checker SaleStateCheckerInterface
	logger  *logan.Entry
	errors  chan error
}

// New is constructor for bearer Service.
func New(config Config, log *logan.Entry, checker SaleStateCheckerInterface) *Service {
	return &Service{
		config:  config,
		checker: checker,
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
	err := s.checkSaleState(ctx)
	if err == nil {
		s.logger.Info("Operation submitted")
		return nil
	}

	if err != errNoSales {
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
