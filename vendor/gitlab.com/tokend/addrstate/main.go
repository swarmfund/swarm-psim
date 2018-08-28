package addrstate

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	horizon "gitlab.com/tokend/horizon-connector"
)

// StateMutator uses to get StateUpdate for specific effects and entryTypes
type StateMutator interface {
	GetStateUpdate(change xdr.LedgerEntryChange) StateUpdate
	GetEffects() []int
	GetEntryTypes() []int
}

// StreamTransactionsV2 streams transactions fetched for specified filters.
type TXStreamerV2 interface {
	StreamTransactionsV2(ctx context.Context, effects, entryTypes []int,
	) (<-chan horizon.TransactionEvent, <-chan error)
}

// Watcher watches what comes from txStreamer and what StateMutators do
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

// New returns new watcher and start streaming transactionsV2
func New(ctx context.Context, log *logan.Entry, mutators []StateMutator, txQ TXStreamerV2) *Watcher {
	ctx, cancel := context.WithCancel(ctx)

	w := &Watcher{
		log:        log.WithField("service", "addrstate"),
		mutators:   mutators,
		txStreamer: txQ,
		ctx:        ctx,

		state:      newState(),
		headUpdate: make(chan struct{}),
	}

	go func() {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.WithRecover(rvr).Error("state watcher panicked")
			}
			cancel()
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

func (w *Watcher) run(ctx context.Context) {
	var entryTypes []int
	var effects []int

	for _, mutator := range w.mutators {
		entryTypes = append(entryTypes, mutator.GetEntryTypes()...)
		effects = append(effects, mutator.GetEffects()...)
	}

	// there is intentionally no recover, it should just die in case of persistent error
	txStream, txStreamErrs := w.txStreamer.StreamTransactionsV2(ctx, effects, entryTypes)

	for {
		select {
		case txEvent := <-txStream:
			if tx := txEvent.Transaction; tx != nil {
				// go through all ledger changes
				for _, change := range tx.Changes {
					// apply all mutators
					ledgerEntryChange, err := convertLedgerEntryChangeV2(change)
					if err != nil {
						w.log.WithError(err).Error("failed to get state update", logan.F{
							"entry_type":     change.EntryType,
							"effect":         change.Effect,
							"transaction_id": tx.ID,
						})
						return
					}
					for _, mutator := range w.mutators {
						stateUpdate := mutator.GetStateUpdate(ledgerEntryChange)
						w.state.Mutate(tx.LedgerCloseTime, stateUpdate)
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
