package addrstate

import (
	"time"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

type StateMutator func(change xdr.LedgerEntryChange) StateUpdate

type TXStreamer interface {
	StreamTransactions(ctx context.Context) (<-chan horizon.TransactionEvent, <-chan error)
}

type Watcher struct {
	log        *logan.Entry
	mutators   []StateMutator
	txStreamer TXStreamer
	ctx        context.Context

	// internal state
	head       time.Time
	headUpdate chan struct{}
	state      *State
}

func New(ctx context.Context, log *logan.Entry, mutators []StateMutator, txQ TXStreamer) *Watcher {
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
	// there is intentionally no recover, it should just die in case of persistent error
	txStream, txStreamErrs := w.txStreamer.StreamTransactions(ctx)

	for {
		select {
		case txEvent := <-txStream:
			if tx := txEvent.Transaction; tx != nil {
				// go through all ledger changes
				for _, change := range tx.LedgerChanges() {
					// apply all mutators
					for _, mutator := range w.mutators {
						w.state.Mutate(tx.CreatedAt, mutator(change))
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
