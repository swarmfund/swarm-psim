package bearer

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/app"
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
		logger:  log.WithField("service", "operation_bearer"),
	}
}

// Run will return closed channel and only when work is finished.
func (s *Service) Run(ctx context.Context) chan error {
	s.logger.Info("Starting.")

	app.RunOverIncrementalTimer(ctx, s.logger, "operation_bearer", s.sendOperations, s.config.Period, 2*s.config.Period)

	errs := make(chan error)
	close(errs)
	return errs
}

// sendOperations is create and submit operations.
func (s *Service) sendOperations(ctx context.Context) error {
	checkSale, err := s.checkSaleState()

	txSuccess, err := s.submitTx(checkSale)
	if txSuccess != nil {
		// txSuccess is not nil only when tx
		// successfully submitted with 200 result code
		s.logger.Debug("Submitted check sale state tx")
		return nil
	}

	serr, ok := err.(horizon.SubmitError)
	if !ok {
		return errors.Wrap(err, "unable to submit tx")
	}

	return errors.Wrap(serr, "tx submission failed", logan.F{
		"response_code":   serr.ResponseCode(),
		"tx_code":         serr.TransactionCode(),
		"operation_codes": serr.OperationCodes(),
	})
}

// submitTx is build transaction, sign and submit it to the Horizon.
func (s *Service) submitTx(ops ...xdr.Operation) (*horizon.TransactionSuccess, error) {
	tb := s.horizon.Transaction(&horizon.TransactionBuilder{
		Source:     s.config.Signer,
		Operations: []xdr.Operation(ops),
	})

	env, err := tb.Sign(s.config.Signer).Marshal64()
	if err != nil {
		return nil, err
	}

	return s.horizon.SubmitTXVerbose(*env)
}
