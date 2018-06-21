package internal

import (
	"time"
)

// BroadcastedEventName type is used for typed constants
type BroadcastedEventName string

// BroadcastedEvent is a structure used to hold data to be sent to analytics services
type BroadcastedEvent struct {
	Account           string
	Name              BroadcastedEventName
	Time              *time.Time
	ActorName         string
	ActorEmail        string
	InvestmentAmount  int64  // Only for "invest-funds event"
	InvestmentCountry string // Only for "invest-funds event"
}

// NewBroadcastedEvent constructs an event filled with data provided by arguments
func NewBroadcastedEvent(Account string, Name BroadcastedEventName, Time *time.Time) *BroadcastedEvent {
	return &BroadcastedEvent{Account, Name, Time, "", "", 0, ""}
}

// WithActor returns a copy of BroadcastedEvent with actor-fields set
func (be *BroadcastedEvent) WithActor(actorName string, actorEmail string) *BroadcastedEvent {
	return &BroadcastedEvent{
		Account:           be.Account,
		Name:              be.Name,
		Time:              be.Time,
		ActorName:         actorName,
		ActorEmail:        actorEmail,
		InvestmentAmount:  be.InvestmentAmount,
		InvestmentCountry: be.InvestmentCountry,
	}
}
