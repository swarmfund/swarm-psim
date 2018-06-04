package listener

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

type AccountProvider interface {
	ByAddress(string) (*horizon.Account, error)
}

type RequestProvider interface {
	GetRequestByID(requestID uint64) (*horizon.Request, error)
}

type TokendHandler struct {
	// TODO by next checkpoint invent better name
	processorByOpType map[xdr.OperationType]Processor

	requestsProvider RequestProvider

	accountsProvider AccountProvider
}

func NewTokendHandler(requestsProvider RequestProvider, accountsProvider AccountProvider) *TokendHandler {
	return &TokendHandler{make(map[xdr.OperationType]Processor), requestsProvider, accountsProvider}
}

func (th *TokendHandler) SetProcessor(opType xdr.OperationType, processor Processor) {
	th.processorByOpType[opType] = processor
}

func (th *TokendHandler) Process(txData <-chan TxData) (<-chan []BroadcastedEvent, error) {
	broadcastedEvents := make(chan []BroadcastedEvent)

	go func() {
		defer func() {
			close(broadcastedEvents)
		}()
		for txDataEntry := range txData {
			txType := txDataEntry.Op.Body.Type
			process := th.processorByOpType[txType]
			if process == nil {
				continue
			}
			broadcastedEvents <- process(txDataEntry)
		}
	}()

	return broadcastedEvents, nil
}
