package pricesetter

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/psim/psim/verification"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

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
