package pricesetter

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/prices/types"
)

type priceFinder interface {
	// TryFind - tries to find most recent PricePoint
	TryFind() (*types.PricePoint, error)
	// RemoveOldPoints - removes points which were created before minAllowedTime
	RemoveOldPoints(minAllowedTime time.Time)
}

type service struct {
	config Config
	log    *logan.Entry

	// TODO Interface
	connector   *horizon.Submitter
	priceFinder priceFinder
	txBuilder   *xdrbuild.Builder
}

func newService(
	config Config,
	log *logan.Entry,
	// TODO Interface
	connector *horizon.Submitter,
	finder priceFinder,
	txBuilder *xdrbuild.Builder) *service {

	return &service{
		config: config,
		log:    log.WithField("service", "price_setter"),

		connector: connector,
		priceFinder: finder,
		txBuilder: txBuilder,
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	app.RunOverIncrementalTimer(ctx, s.log, "rate_sync", s.findAndSubmitPricePoint, 10*time.Second, 5*time.Second)
}

func (s *service) findAndSubmitPricePoint(ctx context.Context) error {
	pointToSubmit, err := s.priceFinder.TryFind()
	if err != nil {
		return errors.Wrap(err, "Failed to find PricePoint to submit")
	}

	if pointToSubmit == nil {
		s.log.Warn("Has not found PricePoint to submit.")
		return nil
	}

	fields := logan.F{
		"price_point": pointToSubmit,
	}

	tx, err := s.txBuilder.Transaction(s.config.Source).Op(xdrbuild.SetAssetPrice{
		BaseAsset:  s.config.BaseAsset,
		QuoteAsset: s.config.QuoteAsset,
		Price:      pointToSubmit.Price,
	}).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal SetAssetPrice TX", fields)
	}

	result := s.connector.Submit(ctx, tx)
	if result.Err != nil {
		return errors.Wrap(result.Err, "Error submitting SetAssetPrice TX to Horizon", fields.Merge(logan.F{
			"submit_result": result,
		}))
	}

	s.log.WithFields(fields).Info("SetAssetPrice TX was submitted to Horizon successfully.")

	s.priceFinder.RemoveOldPoints(pointToSubmit.Time)
	return nil
}
