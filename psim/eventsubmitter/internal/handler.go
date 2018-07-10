package internal

import (
	"context"

	"gitlab.com/tokend/go/xdr"
)

// ProcessedItem holds events generated while processing transactions by an extractor to be emitted by a broadcaster
type ProcessedItem struct {
	BroadcastedEvent *BroadcastedEvent
	Error            error
}

// Handler is responsible for managing which transactions are processed by which processors and for launching them
type Handler interface {
	SetProcessor(opType xdr.OperationType, processor Processor)
	Process(ctx context.Context, extractedItem <-chan ExtractedItem) <-chan ProcessedItem
}

// InvalidProcessedItem returns ProcessedItem only with specified error
func InvalidProcessedItem(err error) *ProcessedItem {
	return &ProcessedItem{BroadcastedEvent: nil, Error: err}
}

// ValidProcessedItem returns ProcessedItem without error, but with specified body
func ValidProcessedItem(event *BroadcastedEvent) *ProcessedItem {
	return &ProcessedItem{BroadcastedEvent: event, Error: nil}
}
