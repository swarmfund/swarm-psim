package state

import "fmt"

// Accounts - represents set of accounts
type Accounts struct {
	// Map of all the account
	Accounts map[AccountID]*Account
}

// NewAccounts - creates new Account instance
func NewAccounts() Accounts {
	return Accounts{
		Accounts: make(map[AccountID]*Account),
	}
}

// Exists - returns true if account already exists in set
func (a *Accounts) Exists(account AccountID) bool {
	_, exists := a.Accounts[account]
	return exists
}

// AddAccount - address account to a. Panics if account was already set
func (a *Accounts) AddAccount(account Account) {
	_, ok := a.Accounts[account.Address]
	if ok {
		panic(fmt.Sprintf("expected account (%s) not to exist", account.Address))
	}

	a.Accounts[account.Address] = &account
}

// GetAccount - returns account by address, panics if account does not exist
func (a *Accounts) GetAccount(address AccountID) *Account {
	account, ok := a.Accounts[address]
	if !ok {
		panic(fmt.Sprintf("expected account (%s) to exist", address))
	}

	return account
}

// GetChildren - provides chanel with accounts which we referred (has parent)
func (a *Accounts) GetChildren() chan *Account {
	accounts := make(chan *Account)
	go func() {
		for _, account := range a.Accounts {
			if account.Parent == "" || account.ShareForReferrer == 0 {
				continue
			}

			accounts <- account
		}

		close(accounts)
	}()
	return accounts
}
