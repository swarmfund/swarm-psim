package listener

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/listener/internal"
	"gitlab.com/tokend/go/xdr"
)

// TokendHandler is a Handler implementation to be used with tokend stuff
type TokendHandler struct {
	processorByOpType map[xdr.OperationType]Processor
}

// NewTokendHandler constructs a TokendHandler without any binded processors
func NewTokendHandler() *TokendHandler {
	return &TokendHandler{make(map[xdr.OperationType]Processor)}
}

// SetProcessor binds a processor to specified opType
func (th TokendHandler) SetProcessor(opType xdr.OperationType, processor Processor) {
	th.processorByOpType[opType] = processor
}

// Process starts processing all op using processors in the map
func (th TokendHandler) Process(ctx context.Context, extractedItems <-chan ExtractedItem) <-chan ProcessedItem {
	broadcastedEvents := make(chan ProcessedItem)

	go func() {
		defer func() {
			close(broadcastedEvents)
		}()

		for extractedItem := range extractedItems {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if extractedItem.Error != nil {
				broadcastedEvents <- *internal.InvalidProcessedItem(errors.Wrap(extractedItem.Error, "received invalid extracted item"))
				continue
			}

			opData := extractedItem.ExtractedOpData
			opType := opData.Op.Body.Type

			process := th.processorByOpType[opType]
			if process == nil {
				continue
			}

			events := process(opData)

			for _, event := range events {
				if event.Error != nil {
					broadcastedEvents <- *internal.InvalidProcessedItem(errors.Wrap(event.Error, "failed to process op data from extracted item"))
					continue
				}

				broadcastedEvents <- *internal.ValidProcessedItem(event.BroadcastedEvent)
			}
		}
	}()

	return broadcastedEvents
}
