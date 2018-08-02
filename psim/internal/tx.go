package internal

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func LedgerChanges(tx *regources.Transaction) []xdr.LedgerEntryChange {
	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXdr, &meta); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal"))
	}
	var result []xdr.LedgerEntryChange
	for _, op := range meta.MustOperations() {
		for _, change := range op.Changes {
			result = append(result, change)
		}
	}
	return result
}

func GroupedLedgerChanges(tx *regources.Transaction) ([][]xdr.LedgerEntryChange, error) {
	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXdr, &meta); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	var result [][]xdr.LedgerEntryChange
	for opIndex, op := range meta.MustOperations() {
		result = append(result, []xdr.LedgerEntryChange{})
		for _, change := range op.Changes {
			result[opIndex] = append(result[opIndex], change)
		}
	}
	return result, nil
}

func SafeEnvelope(tx *regources.Transaction) (*xdr.TransactionEnvelope, error) {
	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshalBase64(tx.EnvelopeXdr, &envelope); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &envelope, nil
}

func Envelope(tx *regources.Transaction) *xdr.TransactionEnvelope {
	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshalBase64(tx.EnvelopeXdr, &envelope); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal"))
	}
	return &envelope
}

func Result(tx *regources.Transaction) (*xdr.TransactionResult, error) {
	var result xdr.TransactionResult
	if err := xdr.SafeUnmarshalBase64(tx.ResultXdr, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal tx result xdr")
	}
	return &result, nil
}
