package kyc

import (
	"gitlab.com/tokend/go/xdrbuild"
	"encoding/json"
	"gitlab.com/tokend/go/xdr"
	"context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/horizon-connector"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type RequestPerformer struct {
	builder *xdrbuild.Builder
	source keypair.Address
	signer keypair.Full
	txSubmitter TXSubmitter
}

func NewRequestPerformer(
	builder *xdrbuild.Builder,
	source keypair.Address,
	signer keypair.Full,
	txSubmitter TXSubmitter) *RequestPerformer {

	return &RequestPerformer {
		builder: builder,
		source: source,
		signer:signer,
		txSubmitter: txSubmitter,
	}
}

func (p *RequestPerformer) Approve(
	ctx context.Context,
	requestID uint64,
	requestHash string,
	tasksToAdd, tasksToRemove uint32,
	extDetails map[string]string) error {

	if extDetails == nil {
		extDetails = make(map[string]string)
	}

	extDetailsBB, err := json.Marshal(extDetails)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal externalDetails", logan.F{"ext_details": extDetails})
	}

	signedEnvelope, err := p.builder.Transaction(p.source).Op(xdrbuild.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.UpdateKYCDetails{
			TasksToAdd:      tasksToAdd,
			TasksToRemove:   tasksToRemove,
			ExternalDetails: string(extDetailsBB),
		},
		Reason: "",
	}).Sign(p.signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal signed Envelope")
	}

	submitResult := p.txSubmitter.Submit(ctx, signedEnvelope)
	if submitResult.Err != nil {
		return errors.Wrap(submitResult.Err, "Error submitting signed Envelope to Horizon", logan.F{
			"submit_result": submitResult,
		})
	}

	return nil
}

func (p *RequestPerformer) Reject(ctx context.Context, requestID uint64, requestHash string, tasksToAdd uint32, extDetails map[string]string, rejectReason string) error {
	if extDetails == nil {
		extDetails = make(map[string]string)
	}

	extDetailsBB, err := json.Marshal(extDetails)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal externalDetails", logan.F{"ext_details": extDetails})
	}

	signedEnvelope, err := p.builder.Transaction(p.source).Op(xdrbuild.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionReject,
		Details: xdrbuild.UpdateKYCDetails{
			TasksToAdd:      tasksToAdd,
			TasksToRemove:   0,
			ExternalDetails: string(extDetailsBB),
		},
		Reason: rejectReason,
	}).Sign(p.signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal signed Envelope")
	}

	submitResult := p.txSubmitter.Submit(ctx, signedEnvelope)
	if submitResult.Err != nil {
		return errors.Wrap(submitResult.Err, "Error submitting signed Envelope to Horizon", logan.F{
			"submit_result": submitResult,
		})
	}

	return nil
}
