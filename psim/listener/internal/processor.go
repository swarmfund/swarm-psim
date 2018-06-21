package internal

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

// OpData holds all the needed stuff used by processors to decide what events to emit
type OpData struct {
	Op              xdr.Operation
	SourceAccount   xdr.AccountId
	OpLedgerChanges []xdr.LedgerEntryChange
	OpResult        xdr.OperationResultTr
	CreatedAt       *time.Time
}

// MaybeBroadcastedEvent can contain BroadcastedEvent OR Error
type MaybeBroadcastedEvent struct {
	BroadcastedEvent *BroadcastedEvent
	Error            error
}

// AppendedBy returns array of the receiver and new event from arguments.
func (mbe *MaybeBroadcastedEvent) AppendedBy(Account string, Name BroadcastedEventName, Time *time.Time) (outputEvents []MaybeBroadcastedEvent) {
	outputEvents = append([]MaybeBroadcastedEvent{*mbe}, MaybeBroadcastedEvent{NewBroadcastedEvent(Account, Name, Time), nil})
	return
}

// Alone returns arrays of one element which is the receiver.
func (mbe *MaybeBroadcastedEvent) Alone() (outputEvents []MaybeBroadcastedEvent) {
	outputEvents = []MaybeBroadcastedEvent{*mbe}
	return
}

// InvalidBroadcastedEvent constructs MaybeBroadcastedEvent with error
func InvalidBroadcastedEvent(err error) *MaybeBroadcastedEvent {
	return &MaybeBroadcastedEvent{nil, err}
}

// ValidBroadcastedEvent constructs MaybeBroadcastedEvent with body
func ValidBroadcastedEvent(Account string, Name BroadcastedEventName, Time *time.Time) *MaybeBroadcastedEvent {
	return &MaybeBroadcastedEvent{NewBroadcastedEvent(Account, Name, Time), nil}
}

// Processor emits events based on data passed in
type Processor func(data OpData) []MaybeBroadcastedEvent
