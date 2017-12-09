package state

import (
	"fmt"
	"time"
)

// State - represents state of the ledger at specific point of time (defined by Ledger - leger sequence number)
type State struct {
	// Ledger Sequence number
	Ledger int64
	// Map of all the account existing in the system
	*Accounts
	// Map Asset-AssetToken
	AssetTokens map[AssetCode]AssetCode

	// Stores map of all special accounts. List of special accounts provided below
	SpecialAccounts *Accounts

	// Master account of the system. Balance of master account stores available for emission coins.
	// It's not allowed to perform transfers to perform transfers to this account.
	// Balances of this account should not participate in any calculations
	MasterAccount AccountID
	// Commission Account - account which stores all the fees paid by users except demurrage fees
	CommissionAccount AccountID
	// Operational Account - account of Digital Gold, should be handled as regular account
	OperationalAccount AccountID
	// StorageFeeAccount - account which stores demurrage fees
	StorageFeeAccount AccountID

	PayoutPeriod *time.Duration
}

// NewState - initializes all the fields of State
func NewState() *State {
	return &State{
		Accounts: &Accounts{
			Accounts: map[AccountID]*Account{},
		},
		SpecialAccounts: &Accounts{
			Accounts: map[AccountID]*Account{},
		},
		AssetTokens: map[AssetCode]AssetCode{},
	}
}

// SetLedger - sets current ledger sequence number
func (state *State) SetLedger(ledger int64) {
	state.Ledger = ledger
}

// GetLedger - returns current ledger sequence number
func (state *State) GetLedger() int64 {
	return state.Ledger
}

// SetToken - sets token for the asset
func (state *State) SetToken(asset, token AssetCode) {
	state.AssetTokens[asset] = token
}

// GetAssetByToken - returns asset by token. Panics if fails to find one
func (state *State) GetAssetByToken(token AssetCode) AssetCode {
	for k, v := range state.AssetTokens {
		if v == token {
			return k
		}
	}
	panic(fmt.Sprintf("expected asset to exist for token (%s)", token))
}

// isToken - returns true, if value is token
func (state *State) isToken(value AssetCode) bool {
	for _, token := range state.AssetTokens {
		if token == value {
			return true
		}
	}
	return false
}

// GetTotalFeesToShare - returns total amount of fees to share in current payout
func (state *State) GetTotalFeesToShare() map[AssetCode]int64 {
	result := map[AssetCode]int64{}
	for _, account := range state.Accounts.Accounts {
		for _, balance := range account.Balances {
			if result[balance.Asset]+balance.GetFeesToShare() < result[balance.Asset] {
				panic(fmt.Sprintf("failed to calculate total fees to share for asset %s - overflow", balance.Asset))
			}
			result[balance.Asset] += balance.GetFeesToShare()
		}
	}
	return result
}

// TokenBalances - returns chanel with tokens balances
func (state *State) TokenBalances() chan *Balance {
	balances := make(chan *Balance)
	go func() {
		for _, account := range state.Accounts.Accounts {
			for _, balance := range account.Balances {
				if state.isToken(balance.Asset) {
					balances <- balance
				}
			}
		}
		close(balances)
	}()
	return balances
}

// GetTotalTokensAmount - returns total amount of each token
func (state *State) GetTotalTokensAmount() map[AssetCode]int64 {
	result := map[AssetCode]int64{}
	for balance := range state.TokenBalances() {
		if result[balance.Asset]+balance.Amount < result[balance.Asset] {
			panic(fmt.Sprintf("failed to calculate total tokens amount for asset %s - overflow", balance.Asset))
		}
		result[balance.Asset] += balance.Amount
	}

	return result
}

// Special Accounts access methods

// GetStorageFeeAccount - returns Storage Fee Account
func (state *State) GetStorageFeeAccount() AccountID {
	return state.StorageFeeAccount
}

// GetOperationalAccount - returns Operational Account
func (state *State) GetOperationalAccount() AccountID {
	return state.OperationalAccount
}

// SetOperationalAccount - sets Operational Account
func (state *State) SetOperationalAccount(accountID AccountID) {
	state.OperationalAccount = accountID
}

// GetMasterAccount - returns Master Account
func (state *State) GetMasterAccount() AccountID {
	return state.MasterAccount
}

// GetCommissionAccount - returns Commission Account
func (state *State) GetCommissionAccount() AccountID {
	return state.CommissionAccount
}

// SetSpecialAccounts - sets all the special accounts
func (state *State) SetSpecialAccounts(master, storage, commission AccountID) {
	state.MasterAccount = master
	state.StorageFeeAccount = storage
	state.CommissionAccount = commission
}

// GetSpecialAccounts - returns special accounts
func (state *State) GetSpecialAccounts() *Accounts {
	return state.SpecialAccounts
}

// GetSpecialAccount - returns special account by its accountID
func (state *State) GetSpecialAccount(accountID AccountID) *Account {
	return state.GetSpecialAccounts().GetAccount(accountID)
}

// IsSpecialAccount - returns true if account requires special handing
func (state *State) IsSpecialAccount(accountID AccountID) bool {
	return state.SpecialAccounts.Exists(accountID)
}

// PayoutCompleted - updates fees after payout for all balances
func (state *State) PayoutCompleted() {
	for _, account := range state.Accounts.Accounts {
		for _, balance := range account.Balances {
			balance.UpdateFeesAfterPayout()
		}
	}
}

// GetMainBalanceForAsset - returns main balance for asset. panics if not found
func (state *State) GetMainBalanceForAsset(accountID AccountID, asset AssetCode) *Balance {
	return state.Accounts.GetAccount(accountID).GetBalanceForAsset(state.GetMasterAccount(), asset)
}

// SetPayoutPeriod - sets payout period
func (state *State) SetPayoutPeriod(d *time.Duration) {
	state.PayoutPeriod = d
}

// GetPayoutPeriod - returns payout period
func (state *State) GetPayoutPeriod() *time.Duration {
	return state.PayoutPeriod
}
