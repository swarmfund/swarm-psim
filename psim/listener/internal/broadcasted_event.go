package internal

import "time"

// BroadcastedEventName type is used for typed constants
type BroadcastedEventName string

// BroadcastedEvent is a structure used to hold data to be sent to analytics services
type BroadcastedEvent struct {
	Account string
	Name    BroadcastedEventName
	Time    time.Time
}

// NewBroadcastedEvent constructs an event filled with data provided by arguments
func NewBroadcastedEvent(Account string, Name BroadcastedEventName, Time time.Time) *BroadcastedEvent {
	return &BroadcastedEvent{Account, Name, Time}
}
