package addrstate

import (
	"time"

	"context"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

const (
	sunAsset = "SUN"
)

type Requester func(ctx context.Context, method, endpoint string, target interface{}) error
type LedgerProvider func(ctx context.Context) <-chan Ledger
type ChangesProvider func(ctx context.Context, ledgerSeq int64) <-chan xdr.LedgerEntryChange
type StateMutator func(change xdr.LedgerEntryChange) StateUpdate

type Watcher struct {
	log       *logan.Entry
	mutator   StateMutator
	ledgers   LedgerProvider
	changes   ChangesProvider
	requester Requester

	head       time.Time
	headUpdate chan struct{}
	state      *State
}

func New(log *logan.Entry, mutator StateMutator, ledgers LedgerProvider, changes ChangesProvider, requester Requester) *Watcher {

	w := &Watcher{
		log:       log.WithField("worker", "address_state_watcher"),
		mutator:   mutator,
		ledgers:   ledgers,
		changes:   changes,
		requester: requester,

		state:      newState(),
		headUpdate: make(chan struct{}),
	}

	go w.run(context.TODO())

	return w
}

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
		case <-ctx.Done():
			return nil
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

func (w *Watcher) BalanceID(ctx context.Context, accountAddress string) (balanceID *string, err error) {
	var response horizon.BalanceIDResponse

	err = w.requester(ctx, "GET", fmt.Sprintf("/accounts/%s/balances", accountAddress), &response)
	if err != nil {
		w.log.WithField("account_address", accountAddress).WithError(err).Error("Failed to request Balances of Account via requester.")
		return nil, err
	}

	if len(response.Balances) == 0 {
		return nil, nil
	}

	for _, b := range response.Balances {
		if b.Asset == sunAsset {
			return &b.BalanceID, nil
		}
	}

	return nil, nil
}

// WatcherState is a connector between LedgerEntryChange and Watcher state for specific consumers
type StateUpdate struct {
	AssetPrice *int64
	Address    *StateAddressUpdate
}

type StateAddressUpdate struct {
	Offchain string
	Tokend   string
}

func (w *Watcher) run(ctx context.Context) {
	for ledger := range w.ledgers(ctx) {
		if ledger.TXCount > 0 {
			for change := range w.changes(ctx, ledger.Sequence) {
				w.state.Mutate(ledger.ClosedAt, w.mutator(change))
			}
		}
		// ledger has been processed bump head
		w.head = ledger.ClosedAt
		w.headUpdate <- struct{}{}
	}
}
