package internal

import (
	"context"
	"time"

	"gitlab.com/tokend/go/xdr"
)

// ExtractedItem holds an error along with data to be passed into the same channel
type ExtractedItem struct {
	ExtractedOpData OpData
	Error           error
}

// Extractor is responsible for gathering all data from transactions it takes from some source to be processed by processors
type Extractor interface {
	Extract(ctx context.Context) <-chan ExtractedItem
}

// InvalidExtractedItem constructs an ExtractedItem with error only
func InvalidExtractedItem(err error) ExtractedItem {
	return ExtractedItem{OpData{}, err}
}

// ValidExtractedItem constructs an ExtractedItem with body and without error
func ValidExtractedItem(currentOp xdr.Operation, sourceAccount xdr.AccountId, opLedgerChanges []xdr.LedgerEntryChange, opResultTr xdr.OperationResultTr, txTime *time.Time) ExtractedItem {
	return ExtractedItem{OpData{Op: currentOp, SourceAccount: sourceAccount, OpLedgerChanges: opLedgerChanges, OpResult: opResultTr, CreatedAt: txTime}, nil}
}
