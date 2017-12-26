package addrstate

import (
	"time"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
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

	go w.run(ctx)

	return w
}

// AddressAt returns the Address of Account, which is coupled with the provided offchain `addr`.
// AddressAt can be blocking, it ensures that returning result
// is based on all the data up to `ts`.
func (w *Watcher) AddressAt(ctx context.Context, ts time.Time, addr string) *string {
	for w.head.Before(ts) {
		select {
		case <-ctx.Done():
			return nil
		case <-w.headUpdate:
			// Make the for check again
			continue
		}
	}

	// Head is not before the `ts` anymore - can respond.

	addr, ok := w.state.addrs[addr]
	if !ok {
		return nil
	}
	return &addr
}

func (w *Watcher) PriceAt(ctx context.Context, ts time.Time) *int64 {
	for w.head.Before(ts) {
		select {
		case <-w.headUpdate:
			// Make the for check again
			continue
		}
	}
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
	// there is intentionally no defer, it should just die in case of persistent error
	transactions := make(chan horizon.Transaction)
	errs := w.txQ.Transactions(transactions)
	for {
		select {
		case tx := <-transactions:
			for _, change := range tx.LedgerChanges() {
				w.state.Mutate(tx.CreatedAt, w.mutator(change))
			}

			w.head = tx.CreatedAt
			w.headUpdate <- struct{}{}
		case err := <-errs:
			w.log.WithError(err).Warn("failed to get transaction")
		}
	}
}
