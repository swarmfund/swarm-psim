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

// CheckSalesStateHelperInterface is an interface allows to perform check sale state operation
// and get all required data for it
type CheckSalesStateHelperInterface interface {
	GetSales() ([]horizon.Sale, error)
	GetHorizonInfo() (info *horizon.Info, err error)
	BuildTx(info *horizon.Info, saleID uint64) (string, error)
	SubmitTx(ctx context.Context, envelope string) horizon.SubmitResult
}

// Service is a main structure for bearer runner,
// implements `app.Service` interface.
type Service struct {
	config Config
	helper CheckSalesStateHelperInterface
	logger *logan.Entry
}

// New is a constructor for bearer Service.
func New(config Config, log *logan.Entry, helper CheckSalesStateHelperInterface) *Service {
	return &Service{
		config: config,
		helper: helper,
		logger: log.WithField("service", conf.ServiceBearer),
	}
}

// Run logs out info about service start and invokes run over incremental timer
func (s *Service) Run(ctx context.Context) {
	s.logger.Info("starting...")

	app.RunOverIncrementalTimer(
		ctx,
		s.logger,
		conf.ServiceBearer,
		s.sendOperations,
		0,
		s.config.AbnormalPeriod)
}

// SendOperations sends operations to Horizon server and gets submission results from it
func (s *Service) sendOperations(ctx context.Context) error {
	err := s.checkSalesState(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot submit checkSalesState operation")
	}

	s.logger.Info("Operation submitted")

	tm := time.NewTimer(s.config.SleepPeriod)
	select {
	case <-ctx.Done():
		return nil
	case <-tm.C:
		return nil
	}
}
