package bearer

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

// Service is a main structure for bearer runner,
// implements `utils.Service` interface.
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

// Run will return closed channel and only when work is finished.
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

// submitOperation is build transaction, sign and submit it to the Horizon.
func (s *Service) submitOperation(op xdr.Operation) error {
	tb := s.horizon.Transaction(&horizon.TransactionBuilder{
		Source:     s.config.Source,
		Operations: []xdr.Operation{op},
	})

	return tb.Sign(s.config.Signer).Submit()
}
