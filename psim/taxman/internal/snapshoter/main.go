package snapshoter

import "gitlab.com/swarmfund/psim/psim/taxman/internal/state"

// OperationSync - helper structure to form core's payment operation
type OperationSync struct {
	Reference            string
	SourceBalanceID      state.BalanceID
	DestinationBalanceID state.BalanceID
	Amount               int64
}

// Snapshot - represents snapshot of the state for which payout will be performed
type Snapshot struct {
	// Ledger is both ledger sequence snapshoter were taken at and unique identifier
	Ledger int64
	// SyncState contains actual stellar operations state that will be packed
	// into transactions during sync step
	SyncState map[string]OperationSync
	// FeesToShareToParent contains total referral payout state for each parent
	FeesToShareToParent map[state.AccountID]map[state.AssetCode]int64
	// FeesToShareToTokenHolders contains tokens shares payout state,
	// `state.AssetCode` key corresponds to asset to be shared
	FeesToShareToTokenHolders map[state.AccountID]map[state.AssetCode]int64
}

// New - creates new instance of snapshoter
func New() *Snapshot {
	return &Snapshot{
		SyncState:                 map[string]OperationSync{},
		FeesToShareToParent:       map[state.AccountID]map[state.AssetCode]int64{},
		FeesToShareToTokenHolders: map[state.AccountID]map[state.AssetCode]int64{},
	}
}

// Snapshots - set of `Snapshot` with Ledger Sequence as key
type Snapshots map[int64]Snapshot

// Add - adds new snapshoter to set of snapshots
func (s Snapshots) Add(snapshot *Snapshot) {
	s[snapshot.Ledger] = *snapshot
}

// Get - returns snapshoter be it's ledger sequence
func (s Snapshots) Get(ledger int64) *Snapshot {
	value, ok := s[ledger]
	if ok {
		return &value
	}
	return nil
}
