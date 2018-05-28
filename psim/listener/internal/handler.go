package internal

import (
	"gitlab.com/tokend/go/xdr"
)

type Handler interface {
	SetProcessor(opType xdr.OperationType, chain Processor)
	Process(txData <-chan TxData) <-chan []BroadcastedEvent
}
