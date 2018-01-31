package bearer

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"fmt"
	"encoding/json"
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

func (s *Service) obtainSales() ([]horizon.Sale, error) {
	respBytes, err := s.horizon.Client().Get(fmt.Sprintf("/core_sales"))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get core sales from Horizon")
	}

	var sales []horizon.Sale
	err = json.Unmarshal(respBytes, &sales)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal Sales from Horizon response", logan.F{
			"horizon_response": string(respBytes),
		})
	}

	return sales, nil
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
