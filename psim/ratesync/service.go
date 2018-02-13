package ratesync

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"time"
	"context"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/tokend/keypair"
)

type priceFinder interface {
	// TryFind - tries to find most recent price point
	TryFind() (*provider.PricePoint, error)
	// RemoveDeprecatedPoints - removes points which were created before minAllowedTime
	RemoveDeprecatedPoints(minAllowedTime time.Time)
}

type service struct {
	baseAsset string
	quoteAsset string


	log *logan.Entry
	source keypair.Address
	signer keypair.Full
	connector *horizon.Submitter
	priceFinder priceFinder
	builder *xdrbuild.Builder
}

// Run is a blocking method, it returns only when ctx closes.
func (s *service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	app.RunOverIncrementalTimer(ctx, s.log, "rate_sync", s.runOnce, 10*time.Second, 5*time.Second)
}

func (s *service) runOnce(ctx context.Context) error {
	pointToSubmit, err := s.priceFinder.TryFind()
	if err != nil {
		return errors.Wrap(err, "failed to find price point to submit")
	}

	if pointToSubmit == nil {
		s.log.Warn("Did not find price point to submit")
		return nil
	}

	tx, err := s.builder.Transaction(s.source).Op(xdrbuild.SetAssetPrice{
		BaseAsset: s.baseAsset,
		QuoteAsset: "SUN",
		Price: pointToSubmit.Price,
	}).Sign(s.signer).Marshal()

	if err != nil {
		return errors.Wrap(err, "failed to create tx")
	}

	result := s.connector.Submit(ctx, tx)
	if result.Err != nil {
		return errors.Wrap(result.Err, "tx was rejected", result.GetLoganFields())
	}

	s.priceFinder.RemoveDeprecatedPoints(pointToSubmit.Time)
	return nil
}


