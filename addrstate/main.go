package addrstate

import (
	"time"
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/regources"
)

type StateMutator interface {
	GetStateUpdate(change regources.LedgerEntryChangeV2) StateUpdate
	GetEffects() []int
	GetEntryTypes() []int
}

type TXStreamerV2 interface {
	StreamTransactionsV2(ctx context.Context, effects, entryTypes []int,
	) (<-chan regources.TransactionV2Event, <-chan error)
}

type Watcher struct {
	log        *logan.Entry
	mutators   []StateMutator
	txStreamer TXStreamerV2
	ctx        context.Context

	// internal state
	head       time.Time
	headUpdate chan struct{}
	state      *State
}

func New(ctx context.Context, log *logan.Entry, mutators []StateMutator, txQ TXStreamerV2) *Watcher {
	w := &Watcher{
		log:        log.WithField("service", "addrstate"),
		mutators:   mutators,
		txStreamer: txQ,

		state:      newState(),
		headUpdate: make(chan struct{}),
	}

	go func() {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.WithRecover(rvr).Error("state watcher panicked")
			}
		}()
		w.run(ctx)
	}()

	return w
}

func (w *Watcher) ensureReached(ctx context.Context, ts time.Time) {
	for w.head.Before(ts) {
		select {
		case <-ctx.Done():
			return
		case <-w.headUpdate:
			// Make the for check again
			continue
		}
	}
}

// WatcherState is a connector between LedgerEntryChange and Watcher state for specific consumers
// if new field added, add case to getStateUpdateTypes method
type StateUpdate struct {
	//AssetPrice      *int64
	ExternalAccount *StateExternalAccountUpdate
	//Address         *StateAddressUpdate
	Balance *StateBalanceUpdate
}

type ExternalAccountBindingState int32

const (
	ExternalAccountBindingStateCreated ExternalAccountBindingState = iota + 1
	ExternalAccountBindingStateDeleted
)

type StateExternalAccountUpdate struct {
	// ExternalType external system accound id type
	ExternalType int32
	// Data external system pool entity data
	Data string
	// Address is a TokenD account address
	Address string
	// State shows current external pool entity binding state
	State ExternalAccountBindingState
}

type StateBalanceUpdate struct {
	Address string
	Balance string
	Asset   string
}

func (w *Watcher) run(ctx context.Context) {
	var entryTypes []int
	var effects []int

	for _, mutator := range w.mutators {
		entryTypes = append(mutator.GetEntryTypes())
		effects = append(mutator.GetEffects())
	}

	// there is intentionally no recover, it should just die in case of persistent error
	txStream, txStreamErrs := w.txStreamer.StreamTransactionsV2(ctx, effects, entryTypes)

	for {
		select {
		case txEvent := <-txStream:
			if tx := txEvent.TransactionV2; tx != nil {
				// go through all ledger changes
				for _, change := range tx.Changes {
					// apply all mutators
					for _, mutator := range w.mutators {
						w.state.Mutate(tx.LedgerCloseTime, mutator.GetStateUpdate(change))
					}
				}
			}

			// if we made it here it's safe to bump head cursor
			w.head = txEvent.Meta.LatestLedger.ClosedAt
			w.headUpdate <- struct{}{}
		case err := <-txStreamErrs:
			w.log.WithError(err).Warn("TXStreamer sent error into channel.")
		}
	}
}
