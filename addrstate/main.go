package addrstate

import (
	"time"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
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
		log:     log.WithField("service", "addrstate"),
		mutator: mutator,
		txQ:     txQ,

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
		case <- ctx.Done():
			return
		case <-w.headUpdate:
			// Make the for check again
			continue
		}
	}
}

func (w *Watcher) AddressAt(ctx context.Context, ts time.Time, offchainAddr string) *string {
	w.ensureReached(ctx, ts)
	if app.IsCanceled(ctx) {
		return nil
	}

	addrI, ok := w.state.addrs.Load(offchainAddr)
	if !ok {
		return nil
	}
	addrValue := addrI.(string)
	return &addrValue
}

func (w *Watcher) Balance(ctx context.Context, address string) *string {
	balanceI, ok := w.state.balances.Load(address)
	if ok {
		balance := balanceI.(string)
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
	balanceI, ok = w.state.balances.Load(address)
	if !ok {
		return nil
	}

	balance := balanceI.(string)
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
