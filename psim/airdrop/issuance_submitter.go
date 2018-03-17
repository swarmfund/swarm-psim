package airdrop

import (
	"context"
	"fmt"

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
	// TODO Remove
	cause           string
	referenceSuffix string

	source keypair.Address
	signer keypair.Full

	builder     *xdrbuild.Builder
	txSubmitter TXSubmitter
}

func NewIssuanceSubmitter(
	asset string,
	// TODO Remove
	cause string,
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
		// TODO Remove
		cause:           cause,
		referenceSuffix: referenceSuffix,

		source:      source,
		signer:      signer,
		builder:     builder,
		txSubmitter: txSubmitter,
	}
}

// Submit returns parameters of the Issuance Operation.
// If reference duplication occurred, Submit returns nil, nil.
// TODO Add details string parameter.
func (s *IssuanceSubmitter) Submit(ctx context.Context, accountAddress, balanceID string, amount uint64) (*issuance.RequestOpt, error) {
	issuanceOpt := issuance.RequestOpt{
		Reference: buildReference(accountAddress, s.referenceSuffix),
		Receiver:  balanceID,
		Asset:     s.asset,
		Amount:    amount,
		Details:   fmt.Sprintf(`{"cause": "%s"}`, s.cause),
	}
	fields := logan.F{
		"issuance_opt": issuanceOpt,
	}

	tx := issuance.CraftIssuanceTX(issuanceOpt, s.builder, s.source, s.signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal TX into Envelope", fields)
	}

	ok, err := issuance.SubmitEnvelope(ctx, envelope, s.txSubmitter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to submit IssuanceRequest TX Envelope to Horizon", fields)
	}

	if ok {
		return &issuanceOpt, nil
	} else {
		// Reference duplication
		return nil, nil
	}
}

func buildReference(accountAddress, suffix string) string {
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
