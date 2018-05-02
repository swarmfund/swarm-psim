package pricesetter

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/prices/types"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
)

var (
	ErrNoVerifierServices = errors.New("No Deposit Verify services were found.")
)

type priceFinder interface {
	// TryFind - tries to find most recent PricePoint
	TryFind() (*types.PricePoint, error)
	// RemoveOldPoints - removes points which were created before minAllowedTime
	RemoveOldPoints(minAllowedTime time.Time)
}

// Discovery must be implemented by Discovery(Consul) client to pass into Service constructor.
type Discovery interface {
	DiscoverService(service string) ([]discovery.ServiceEntry, error)
}

type service struct {
	config Config
	log    *logan.Entry

	// TODO Interface
	connector   *horizon.Submitter
	priceFinder priceFinder
	txBuilder   *xdrbuild.Builder
	discovery   Discovery
}

func newService(
	config Config,
	log *logan.Entry,
	// TODO Interface
	connector *horizon.Submitter,
	finder priceFinder,
	txBuilder *xdrbuild.Builder,
	discovery Discovery) *service {

	return &service{
		config: config,
		log: log.WithField("service", "price_setter").WithField("base_asset", config.BaseAsset).
			WithField("quote_asset", config.QuoteAsset),

		connector:   connector,
		priceFinder: finder,
		txBuilder:   txBuilder,
		discovery:   discovery,
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	running.WithBackOff(ctx, s.log, "price_setter", s.findAndProcessPricePoint, 10*time.Second, 5*time.Second, 5*time.Minute)
}

func (s *service) findAndProcessPricePoint(ctx context.Context) error {
	pointToSubmit, findErr := s.priceFinder.TryFind()
	if findErr != nil {
		s.log.WithError(findErr).Warn("Has not found PricePoint to submit.")
		return nil
	}

	fields := logan.F{
		"price_point": pointToSubmit,
	}
	s.log.WithFields(fields).Info("Found PricePoint meeting restrictions.")

	envelope, err := s.txBuilder.Transaction(s.config.Source).Op(xdrbuild.SetAssetPrice{
		BaseAsset:  s.config.BaseAsset,
		QuoteAsset: s.config.QuoteAsset,
		Price:      pointToSubmit.Price,
	}).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal SetAssetPrice TX", fields)
	}

	verifiedEnvelope, err := s.verifyEnvelope(envelope, pointToSubmit.Price)
	if err != nil {
		return errors.Wrap(err, "Failed to verify Envelope", fields)
	}

	result := s.connector.Submit(ctx, verifiedEnvelope)
	if result.Err != nil {
		return errors.Wrap(result.Err, "Error submitting SetAssetPrice TX to Horizon", fields.Merge(logan.F{
			"submit_result": result,
		}))
	}

	s.log.WithFields(fields).Info("SetAssetPrice TX was submitted to Horizon successfully.")

	s.priceFinder.RemoveOldPoints(pointToSubmit.Time)
	return nil
}
