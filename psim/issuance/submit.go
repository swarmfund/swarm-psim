package issuance

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

// SubmitEnvelope is a helper to handle OpReferenceDuplication.
// If reference duplication happens - false and nil error will be returned,
// in case all other errors - non-nil error will be returned.
func SubmitEnvelope(ctx context.Context, envelope string, submitter TXSubmitter) (bool, error) {
	result := submitter.Submit(ctx, envelope)
	if result.Err != nil {
		if len(result.OpCodes) == 1 && result.OpCodes[0] == "op_reference_duplication" {
			// Deposit duplication - we already processed this deposit - just ignoring it.
			return false, nil
		}

		return false, errors.Wrap(result.Err, "Horizon SubmitResult has error", logan.F{
			"submit_result": result,
		})
	}

	return true, nil
}
