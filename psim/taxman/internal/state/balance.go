package state

// Balance - represent core balance entry
type Balance struct {
	// Account id of the owner
	Account AccountID
	// Asset of the account
	Asset AssetCode
	// Address of the balance
	Address BalanceID
	// Balance amount
	Amount int64
	// ID of the exchange if which account was opened
	ExchangeID AccountID

	// FeesPaid - fees paid since last demurrage
	FeesPaid int64
	// FeesSharedButNotCleared stores FeesPaid, which are already shared by payout, but still stored in FeesPaid
	FeesSharedButNotCleared int64
	// FeesClearedButNotShared stores sum of FeesPaid for which we have not performed payout, but they were set to 0 by demurrage
	FeesClearedButNotShared int64
}

// SetFeesPaid - sets fees paid for specified balance.
func (b *Balance) SetFeesPaid(feesPaid int64) {
	// smells like demurrage
	if feesPaid == 0 {
		b.updateFeesAfterDemurrage()
		return
	}

	b.FeesPaid = feesPaid
}

// updateFeesAfterDemurrage - updates fees after demurrage
func (b *Balance) updateFeesAfterDemurrage() {
	// demurrage sets fees paid to 0, as we have not performed payout for this fees yet, we should store them
	b.FeesClearedButNotShared += b.FeesPaid
	b.FeesPaid = 0
	b.FeesSharedButNotCleared = 0
}

// UpdateFeesAfterPayout - updates fees after payout.
func (b *Balance) UpdateFeesAfterPayout() {
	b.FeesSharedButNotCleared = b.FeesPaid
	b.FeesClearedButNotShared = 0
}

// GetFeesToShare - returns amount of fees to share for current payout
func (b *Balance) GetFeesToShare() int64 {
	return b.FeesPaid + b.FeesClearedButNotShared - b.FeesSharedButNotCleared
}
