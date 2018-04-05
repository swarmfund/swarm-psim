package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/issuance"
	"gitlab.com/tokend/keypair"
)

const (
	maxReferenceLen = 64
	accountIDLen    = 56
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type IssuanceSubmitter struct {
	asset           string
	referenceSuffix string

	source keypair.Address
	signer keypair.Full

	builder     *xdrbuild.Builder
	txSubmitter TXSubmitter
}

func NewIssuanceSubmitter(
	asset string,
	referenceSuffix string,
	source keypair.Address,
	signer keypair.Full,
	builder *xdrbuild.Builder,
	txSubmitter TXSubmitter) *IssuanceSubmitter {

	if len(referenceSuffix) != maxReferenceLen-accountIDLen {
		panic(errors.Errorf("ReferenceSuffix length must be exactly %d.", maxReferenceLen-accountIDLen))
		return nil
	}

	return &IssuanceSubmitter{
		asset:           asset,
		referenceSuffix: referenceSuffix,

		source:      source,
		signer:      signer,
		builder:     builder,
		txSubmitter: txSubmitter,
	}
}

// Submit returns parameters of the Issuance Operation.
// If reference duplication occurred, Submit returns nil, nil.
func (s *IssuanceSubmitter) Submit(ctx context.Context, accountAddress, balanceID string, amount uint64, opDetails string) (*issuance.RequestOpt, bool, error) {
	issuanceOpt := issuance.RequestOpt{
		Reference: BuildReference(accountAddress, s.referenceSuffix),
		Receiver:  balanceID,
		Asset:     s.asset,
		Amount:    amount,
		Details:   opDetails,
	}
	fields := logan.F{
		"issuance_opt": issuanceOpt,
	}

	tx := issuance.CraftIssuanceTX(issuanceOpt, s.builder, s.source, s.signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return nil, false, errors.Wrap(err, "Failed to marshal TX into Envelope", fields)
	}

	ok, err := issuance.SubmitEnvelope(ctx, envelope, s.txSubmitter)
	if err != nil {
		return nil, false, errors.Wrap(err, "Failed to submit IssuanceRequest TX Envelope to Horizon", fields)
	}

	return &issuanceOpt, ok, nil
}

func BuildReference(accountAddress, suffix string) string {
	result := accountAddress + suffix

	// Just in case.
	if len(result) > maxReferenceLen {
		result = result[len(result)-maxReferenceLen:]
	}

	// Just in case.
	if len(result) < maxReferenceLen {
		filler := "----------------------------------------------------------------" // len = 64
		result = result + filler[:maxReferenceLen-len(result)]
	}

	return result
}
