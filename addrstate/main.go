package addrstate

import (
	"time"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
)

type Watcher struct {
	log     *logan.Entry
	mutator StateMutator
	txQ     TransactionQ
	ctx     context.Context

	// internal state
	head       time.Time
	headUpdate chan struct{}
	state      *State
}

func New(ctx context.Context, log *logan.Entry, mutator StateMutator, txQ TransactionQ) *Watcher {
	w := &Watcher{
		log:     log,
		mutator: mutator,
		txQ:     txQ,

		state:      newState(),
		headUpdate: make(chan struct{}),
	}

	go func() {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.WithError(errors.FromPanic(rvr)).Error("state watcher panicked")
			}
		}()
		w.run(ctx)
	}()

	return w
}

// ensureReached will block until state head reached provided ts
func (w *Watcher) ensureReached(ts time.Time) {
	for w.head.Before(ts) {
		select {
		case <-w.headUpdate:
			// Make the for check again
			continue
		}
	}
}

func (w *Watcher) AddressAt(ctx context.Context, ts time.Time, addr string) *string {
	w.ensureReached(ts)

	addr, ok := w.state.addrs[addr]
	if !ok {
		return nil
	}
	return &addr
}

func (w *Watcher) PriceAt(ctx context.Context, ts time.Time) *int64 {
	w.ensureReached(ts)

	for _, price := range w.state.prices {
		if ts.After(price.UpdatedAt) {
			return &price.Value
		}
	}
	return nil
}

func (w *Watcher) Balance(ctx context.Context, address string) *string {
	balance, ok := w.state.balances[address]
	if ok {
		return &balance
	}
	// if we don't have balance already, let's wait for latest ledger
	now := time.Now()
	for w.head.Before(now) {
		select {
		case <-w.headUpdate:
			continue
		}
	}
	// now check again
	balance, ok = w.state.balances[address]
	if !ok {
		return nil
	}
	return &balance
}

// WatcherState is a connector between LedgerEntryChange and Watcher state for specific consumers
type StateUpdate struct {
	AssetPrice *int64
	Address    *StateAddressUpdate
	Balance    *StateBalanceUpdate
}

type StateAddressUpdate struct {
	Offchain string
	Tokend   string
}

type StateBalanceUpdate struct {
	Address string
	Balance string
}

func (w *Watcher) run(ctx context.Context) {
	// there is intentionally no recover, it should just die in case of persistent error
	events := make(chan horizon.TransactionEvent)
	errs := w.txQ.Transactions(events)
	for {
		select {
		case event := <-events:
			if tx := event.Transaction; tx != nil {
				for _, change := range tx.LedgerChanges() {
					w.state.Mutate(tx.CreatedAt, w.mutator(change))
				}
			}
			w.head = event.Meta.LatestLedger.ClosedAt
			w.headUpdate <- struct{}{}
		case err := <-errs:
			w.log.WithError(err).Warn("failed to get transaction")
		}
	}
}
