package internal

import "time"

type BroadcastedEventName string

type BroadcastedEvent struct {
	Account string
	Name    BroadcastedEventName
	Time    *time.Time
}

func NewBroadcastedEvent(Account string, Name BroadcastedEventName, Time *time.Time) *BroadcastedEvent {
	return &BroadcastedEvent{Account, Name, Time}
}

// AppendedBy returns array of the receiver and new event from arguments.
func (oe *BroadcastedEvent) AppendedBy(Account string, Name BroadcastedEventName, Time *time.Time) (outputEvents []BroadcastedEvent) {
	outputEvents = append([]BroadcastedEvent{*oe}, *NewBroadcastedEvent(Account, Name, Time))
	return
}

// Alone returns arrays of one element which is the receiver.
func (oe *BroadcastedEvent) Alone() (outputEvents []BroadcastedEvent) {
	outputEvents = []BroadcastedEvent{*oe}
	return
}
