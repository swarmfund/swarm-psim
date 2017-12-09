package state

import (
	"fmt"
)

// Account - represent core's account entry
type Account struct {
	// Address of the account
	Address AccountID
	// Balances of the account
	Balances []*Balance
	// Referrer of the account
	Parent AccountID
	// Fee share for referrer
	ShareForReferrer int64
}

// GetBalanceForAsset - returns balance for asset. If balance not found - returns nil
func (a *Account) GetBalanceForAsset(exchangeID AccountID, asset AssetCode) *Balance {
	for _, balance := range a.Balances {
		if balance.Asset == asset && balance.ExchangeID == exchangeID {
			return balance
		}
	}

	return nil
}

// MustGetBalanceForBalanceID - returns balance for balance ID. Panics if balance does not exist
func (a *Account) MustGetBalanceForBalanceID(balanceID BalanceID) *Balance {
	balance := a.tryGetBalanceForBalanceID(balanceID)
	if balance != nil {
		return balance
	}

	panic(fmt.Sprintf("Invalid state: account %s does not have balance with ID %s", a.Address, balanceID))
}

func (a *Account) tryGetBalanceForBalanceID(balanceID BalanceID) *Balance {
	for _, balance := range a.Balances {
		if balance.Address == balanceID {
			return balance
		}
	}

	return nil
}

// AddBalance - adds balance to account. If balance with such balanceID already exist returns error
func (a *Account) AddBalance(balance Balance) error {
	storedBalance := a.tryGetBalanceForBalanceID(balance.Address)
	if storedBalance != nil {
		return fmt.Errorf("invalid state: %s already have balance with ID %s", a.Address, balance.Address)
	}

	a.Balances = append(a.Balances, &balance)
	return nil
}
