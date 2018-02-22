package issuance

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

var (
	OpCodeReferenceDuplication = "op_reference_duplication"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

func SubmitEnvelope(ctx context.Context, envelope xdr.TransactionEnvelope, submitter TXSubmitter) (bool, error) {
	envelopeBase64, err := xdr.MarshalBase64(envelope)
	if err != nil {
		return false, errors.Wrap(err, "Failed to marshal fully signed Envelope")
	}

	result := submitter.Submit(ctx, envelopeBase64)
	if result.Err != nil {
		if len(result.OpCodes) == 1 && result.OpCodes[0] == OpCodeReferenceDuplication {
			// Deposit duplication - we already processed this deposit - just ignoring it.
			return false, nil
		}

		return false, errors.Wrap(result.Err, "Horizon SubmitResult has error", logan.F{
			"submit_result": result,
		})
	}

	return true, nil
}
