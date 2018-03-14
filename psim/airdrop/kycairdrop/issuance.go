package kycairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/issuance"
	"fmt"
	"gitlab.com/swarmfund/psim/psim/airdrop"
)

// SubmitIssuance returns parameters of the Issuance Operation.
// If reference duplication occurred, SubmitIssuance returns nil, nil.
func (s *Service) submitIssuance(ctx context.Context, accountAddress, balanceID string) (*issuance.RequestOpt, error) {
	issuanceOpt := issuance.RequestOpt{
		Reference: buildReference(accountAddress),
		Receiver:  balanceID,
		Asset:     s.config.Asset,
		Amount:    s.config.Amount,
		Details:   fmt.Sprintf(`{"cause": "%s"}`, airdrop.KYCIssuanceCause),
	}
	fields := logan.F{
		"issuance_opt": issuanceOpt,
	}

	tx := issuance.CraftIssuanceTX(issuanceOpt, s.builder, s.config.Source, s.config.Signer)

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

func buildReference(accountAddress string) string {
	const maxReferenceLen = 64

	result := accountAddress + airdrop.KYCReferenceSuffix

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
