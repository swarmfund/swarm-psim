package pricesetterveri

import (
	"time"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/psim/psim/app"
)

var (
	pointsCleaningPeriod = 5 * time.Minute
)

type priceFinder interface {
	VerifyPrice(price int64) error
	// RemoveOldPoints - removes points which were created before minAllowedTime
	RemoveOldPoints(minAllowedTime time.Time)
}

type Verifier struct {
	log         *logan.Entry
	config      Config
	priceFinder priceFinder
}

func newVerifier(
	serviceName string,
	log *logan.Entry,
	config Config,
	priceFinder priceFinder) *Verifier {

	return &Verifier{
		log:         log.WithField("service", serviceName),
		config:      config,
		priceFinder: priceFinder,
	}
}

func (v *Verifier) Run(ctx context.Context) {
	app.RunOverIncrementalTimer(ctx, v.log, "price_points_cleaner", v.cleanPricePoints, pointsCleaningPeriod, 30*time.Second)
}

// CleanPricePoints always returns nil. Returning error - is just to fit the
// signature of a function needed for RunOverIncrementalTimer.
func (v *Verifier) cleanPricePoints(ctx context.Context) error {
	v.priceFinder.RemoveOldPoints(time.Now().Add(-pointsCleaningPeriod))
	return nil
}

func (v *Verifier) GetOperationType() xdr.OperationType {
	return xdr.OperationTypeManageAssetPair
}

func (v *Verifier) VerifyOperation(envelope xdr.TransactionEnvelope) (verifyErr, err error) {
	op := envelope.Tx.Operations[0].Body.ManageAssetPairOp

	if op == nil {
		return errors.Errorf("ManageAssetPairOp is nil."), nil
	}

	verifyErr, err = v.verifyManageAssetPairOp(*op)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to validate Issuance Op")
	}

	return verifyErr, nil
}

func (v *Verifier) verifyManageAssetPairOp(op xdr.ManageAssetPairOp) (verifyErr, err error) {
	if string(op.Base) != v.config.BaseAsset {
		return errors.Errorf("Invalid BaseAsset, expected (%s), got (%s)", v.config.BaseAsset, op.Base), nil
	}
	if string(op.Quote) != v.config.QuoteAsset {
		return errors.Errorf("Invalid QuoteAsset, expected (%s), got (%s)", v.config.QuoteAsset, op.Quote), nil
	}

	if op.Action != xdr.ManageAssetPairActionUpdatePrice {
		return errors.Errorf("Invalid Operation Action, expected UpdatePrice(%d), got (%d)", xdr.ManageAssetPairActionUpdatePrice, op.Action), nil
	}

	priceVerifyErr := v.priceFinder.VerifyPrice(int64(op.PhysicalPrice))
	if priceVerifyErr != nil {
		return errors.Wrap(priceVerifyErr, "Price is invalid"), nil
	}

	return nil, nil
}
