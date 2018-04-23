package bearer

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

type SalesQ interface {
	Sales() ([]horizon.Sale, error)
}

// CheckSalesStateHelperInterface is an interface allows to perform check sale state operation
// and get all required data for it
type CheckSalesStateHelperInterface interface {
	SalesQ
	CloseSale(id uint64) (bool, error)
	//GetHorizonInfo() (info *horizon.Info, err error)
	//BuildTx(info *horizon.Info, saleID uint64) (string, error)
	//SubmitTx(ctx context.Context, envelope string) horizon.SubmitResult
}

// Service is a main structure for bearer runner,
// implements `app.Service` interface.
type Service struct {
	config Config
	helper CheckSalesStateHelperInterface
	logger *logan.Entry
	ticker *time.Ticker
}

// New is a constructor for bearer Service.
func New(config Config, log *logan.Entry, helper CheckSalesStateHelperInterface) *Service {
	return &Service{
		config: config,
		helper: helper,
		logger: log.WithField("service", conf.ServiceBearer),
		ticker: time.NewTicker(config.SleepPeriod),
	}
}

// Run logs out info about service start and invokes run over incremental timer
func (s *Service) Run(ctx context.Context) {
	s.logger.Info("starting...")

	app.RunOverIncrementalTimer(
		ctx,
		s.logger,
		conf.ServiceBearer,
		s.worker,
		0,
		s.config.AbnormalPeriod)
}

// SendOperations sends operations to Horizon server and gets submission results from it
func (s *Service) worker(ctx context.Context) error {
	if err := s.checkSalesState(ctx); err != nil {
		return errors.Wrap(err, "cannot submit checkSalesState operation")
	}

	s.logger.Info("successful iteration")

	select {
	case <-ctx.Done():
		s.logger.Info("bye-bye")
	case <-s.ticker.C:
	}

	return nil
}

func (s *Service) checkSalesState(ctx context.Context) error {
	sales, err := s.helper.Sales()
	if err != nil {
		return errors.Wrap(err, "failed to get sales")
	}

	for _, sale := range sales {
		fields := logan.F{
			"sale_id": sale.ID,
		}
		closed, err := s.helper.CloseSale(sale.ID)
		if err != nil {
			return errors.Wrap(err, "failed to close sale", fields)
		}
		if closed {
			s.logger.WithFields(fields).Info("sale closed")
		} else {
			s.logger.WithFields(fields).Info("sale not ready yet")
		}
	}

	return nil
}
