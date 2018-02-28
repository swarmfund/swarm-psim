package pricesetter

import (
	"context"
	"time"

	"fmt"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/prices/types"
	"gitlab.com/swarmfund/psim/psim/verification"
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
		log:    log.WithField("service", "price_setter"),

		connector:   connector,
		priceFinder: finder,
		txBuilder:   txBuilder,
		discovery:   discovery,
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	app.RunOverIncrementalTimer(ctx, s.log, "price_setter", s.findAndProcessPricePoint, 10*time.Second, 5*time.Second)
}

func (s *service) findAndProcessPricePoint(ctx context.Context) error {
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
	s.log.WithFields(fields).Info("Found Point meeting restrictions.")

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

func (s *service) verifyEnvelope(envelope string, price int64) (string, error) {
	readyEnvelope, err := s.sendToVerifier(envelope)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send to Verifier of verification unsuccessful")
	}

	checkErr := s.checkVerifiedEnvelope(*readyEnvelope, price)
	if checkErr != nil {
		return "", errors.Wrap(err, "Fully signed Envelope from Verifier is invalid")
	}

	envelopeBase64, err := xdr.MarshalBase64(*readyEnvelope)
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal fully signed Envelope")
	}

	return envelopeBase64, nil
}

func (s *service) sendToVerifier(envelope string) (fullySignedTXEnvelope *xdr.TransactionEnvelope, err error) {
	url, err := s.getVerifierURL()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get URL of Verify")
	}

	responseEnvelope, err := verification.Verify(url, envelope)
	if err != nil {
		return nil, errors.Wrap(err, "Verification unsuccessful", logan.F{"verifier_url": url})
	}

	return responseEnvelope, nil
}

func (s *service) getVerifierURL() (string, error) {
	services, err := s.discovery.DiscoverService(s.config.VerifierServiceName)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to discover %s service.", s.config.VerifierServiceName))
	}
	if len(services) == 0 {
		return "", ErrNoVerifierServices
	}

	return services[0].Address, nil
}

func (s *service) checkVerifiedEnvelope(envelope xdr.TransactionEnvelope, price int64) (checkErr error) {
	if len(envelope.Tx.Operations) != 1 {
		return errors.New("Must be exactly 1 Operation.")
	}

	opBody := envelope.Tx.Operations[0].Body

	if opBody.Type != xdr.OperationTypeManageAssetPair {
		return errors.Errorf("Expected OperationType to be ManageAssetPair(%d), but got (%d).",
			xdr.OperationTypeManageAssetPair, opBody.Type)
	}

	op := envelope.Tx.Operations[0].Body.ManageAssetPairOp

	if op == nil {
		return errors.Errorf("ManageAssetPairOp is nil.")
	}

	if string(op.Base) != s.config.BaseAsset {
		return errors.Errorf("Invalid BaseAsset, expected (%s), got (%s)", s.config.BaseAsset, op.Base)
	}
	if string(op.Quote) != s.config.QuoteAsset {
		return errors.Errorf("Invalid QuoteAsset, expected (%s), got (%s)", s.config.QuoteAsset, op.Quote)
	}

	if op.Action != xdr.ManageAssetPairActionUpdatePrice {
		return errors.Errorf("Invalid Operation Action, expected UpdatePrice(%d), got (%d)", xdr.ManageAssetPairActionUpdatePrice, op.Action)
	}

	if int64(op.PhysicalPrice) != price {
		return errors.Errorf("Price is invalid, expected (%d), got (%d).", price, op.PhysicalPrice)
	}

	return nil
}
