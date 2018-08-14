package internal

import (
	"time"
)

// BroadcastedEventName type is used for typed constants
type BroadcastedEventName string

// BroadcastedEvent is a structure used to hold data to be sent to analytics services
type BroadcastedEvent struct {
	Account          string
	Name             BroadcastedEventName
	Time             time.Time
	ActorName        string
	ActorEmail       string
	InvestmentAmount int64 // Only for "invest-funds event"
	DepositAmount    int64
	DepositCurrency  string
	Referral         string
	Country          string
}

// NewBroadcastedEvent constructs an event filled with data provided by arguments
func NewBroadcastedEvent(Account string, Name BroadcastedEventName, Time time.Time) *BroadcastedEvent {
	return &BroadcastedEvent{
		Account: Account,
		Name:    Name,
		Time:    Time,
	}
}

// WithActor returns a copy of BroadcastedEvent with actor-fields set
func (b BroadcastedEvent) WithActor(actorName string, actorEmail string) *BroadcastedEvent {
	newB := b

	newB.ActorName = actorName
	newB.ActorEmail = actorEmail

	return &newB
}
