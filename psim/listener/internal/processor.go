package internal

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

type TxData struct {
	Op              xdr.Operation
	SourceAccount   xdr.AccountId
	OpLedgerChanges []xdr.LedgerEntryChange
	OpResult        xdr.OperationResultTr
	CreatedAt       *time.Time
}

type Processor func(d TxData) []BroadcastedEvent
