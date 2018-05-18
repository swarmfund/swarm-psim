package kyc

import (
	"context"
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

const (
	RejectorExtDetailsKey string = "rejector"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type RequestPerformer struct {
	builder     *xdrbuild.Builder
	source      keypair.Address
	signer      keypair.Full
	txSubmitter TXSubmitter
}

func NewRequestPerformer(
	builder *xdrbuild.Builder,
	source keypair.Address,
	signer keypair.Full,
	txSubmitter TXSubmitter) *RequestPerformer {

	return &RequestPerformer{
		builder:     builder,
		source:      source,
		signer:      signer,
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

func (p *RequestPerformer) Reject(
	ctx context.Context,
	requestID uint64,
	requestHash string,
	tasksToAdd uint32,
	extDetails map[string]string,
	rejectReason, rejector string) error {

	if extDetails == nil {
		extDetails = make(map[string]string)
	}
	extDetails[RejectorExtDetailsKey] = rejector

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
