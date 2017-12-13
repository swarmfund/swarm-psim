package addrstate

import (
	"time"

	"gitlab.com/swarmfund/go/xdr"
)

type Requester func(method, endpoint string, target interface{}) error
type LedgerProvider func() <-chan Ledger
type ChangesProvider func(ledgerSeq int64) <-chan xdr.LedgerEntryChange
type StateMutator func(change xdr.LedgerEntryChange) StateUpdate

type Watcher struct {
	ledgers    LedgerProvider
	changes    ChangesProvider
	mutator    StateMutator
	head       time.Time
	headUpdate chan struct{}
	state      *State
}

func New(mutator StateMutator, ledgers LedgerProvider, changes ChangesProvider) *Watcher {
	w := &Watcher{
		mutator:    mutator,
		ledgers:    ledgers,
		changes:    changes,
		state:      newState(),
		headUpdate: make(chan struct{}),
	}
	go w.run()
	return w
}

func (w *Watcher) AddressAt(ts time.Time, addr string) *string {
	for w.head.Before(ts) {
		<-w.headUpdate
	}
	addr, ok := w.state.addrs[addr]
	if !ok {
		return nil
	}
	return &addr
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

func (w *Watcher) run() {
	for ledger := range w.ledgers() {
		if ledger.TXCount > 0 {
			for change := range w.changes(ledger.Sequence) {
				w.state.Mutate(ledger.ClosedAt, w.mutator(change))
			}
		}
		// ledger has been processed bump head
		w.head = ledger.ClosedAt
		w.headUpdate <- struct{}{}
	}
}
