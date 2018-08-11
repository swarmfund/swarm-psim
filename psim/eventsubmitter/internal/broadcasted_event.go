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
	return &BroadcastedEvent{Account, Name, Time, "", "", 0, 0, "", "", ""}
}

// WithActor returns a copy of BroadcastedEvent with actor-fields set
func (be *BroadcastedEvent) WithActor(actorName string, actorEmail string) *BroadcastedEvent {
	return &BroadcastedEvent{
		Account:          be.Account,
		Name:             be.Name,
		Time:             be.Time,
		ActorName:        actorName,
		ActorEmail:       actorEmail,
		InvestmentAmount: be.InvestmentAmount,
		DepositAmount:    be.DepositAmount,
		DepositCurrency:  be.DepositCurrency,
		Referral:         be.Referral,
		Country:          be.Country,
	}
}

// WithActor returns a copy of BroadcastedEvent with actor-fields set
func (be *BroadcastedEvent) WithDeposit(amount int64, currency string) *BroadcastedEvent {
	return &BroadcastedEvent{
		Account:          be.Account,
		Name:             be.Name,
		Time:             be.Time,
		ActorName:        be.ActorName,
		ActorEmail:       be.ActorEmail,
		InvestmentAmount: be.InvestmentAmount,
		DepositAmount:    amount,
		DepositCurrency:  currency,
		Referral:         be.Referral,
		Country:          be.Country,
	}
}
