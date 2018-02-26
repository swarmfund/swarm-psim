package pricesetter

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

type priceFinder interface {
	// TryFind - tries to find most recent PricePoint
	TryFind() (*providers.PricePoint, error)
	// RemoveOldPoints - removes points which were created before minAllowedTime
	RemoveOldPoints(minAllowedTime time.Time)
}

type service struct {
	baseAsset  string
	quoteAsset string

	log *logan.Entry

	source keypair.Address
	signer keypair.Full

	connector   *horizon.Submitter
	priceFinder priceFinder
	txBuilder   *xdrbuild.Builder
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

	tx, err := s.txBuilder.Transaction(s.source).Op(xdrbuild.SetAssetPrice{
		BaseAsset:  s.baseAsset,
		QuoteAsset: s.quoteAsset,
		Price:      pointToSubmit.Price,
	}).Sign(s.signer).Marshal()
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
