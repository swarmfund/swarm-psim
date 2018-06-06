package resources

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdr"
)

type Transaction struct {
	CreatedAt     time.Time `json:"created_at"`
	PagingToken   string    `json:"paging_token"`
	ResultMetaXDR string    `json:"result_meta_xdr"`
	EnvelopeXDR   string    `json:"envelope_xdr"`
	ResultXDR     string    `json:"result_xdr"`
}

// returns flat array with all the ledger changes
func (tx *Transaction) LedgerChanges() []xdr.LedgerEntryChange {
	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXDR, &meta); err != nil {
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

// returns array of ledger changes for every operation in tx
func (tx *Transaction) GroupedLedgerChanges() ([][]xdr.LedgerEntryChange, error) {
	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXDR, &meta); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	var result [][]xdr.LedgerEntryChange
	for opIndex, op := range meta.MustOperations() {
		result = append(result, []xdr.LedgerEntryChange {})
		for _, change := range op.Changes {
			result[opIndex] = append(result[opIndex], change)
		}
	}
	return result, nil
}

func (tx *Transaction) SafeEnvelope() (*xdr.TransactionEnvelope, error) {
	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshalBase64(tx.EnvelopeXDR, &envelope); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &envelope, nil
}

func (tx *Transaction) Envelope() *xdr.TransactionEnvelope {
	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshalBase64(tx.EnvelopeXDR, &envelope); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal"))
	}
	return &envelope
}

func (tx *Transaction) Result() (*xdr.TransactionResult, error) {
	var result xdr.TransactionResult
	if err := xdr.SafeUnmarshalBase64(tx.ResultXDR, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal tx result xdr")
	}
	return &result, nil
}

func (tx Transaction) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"created_at": tx.CreatedAt,
		"paging_token": tx.PagingToken,
	}
}
