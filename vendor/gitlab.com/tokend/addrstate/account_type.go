package addrstate

import "gitlab.com/tokend/go/xdr"

type AccountTypeMutator struct{}

func (m AccountTypeMutator) GetEffects() []int {
	return []int{int(xdr.LedgerEntryChangeTypeUpdated), int(xdr.LedgerEntryChangeTypeCreated)}
}

func (m AccountTypeMutator) GetEntryTypes() []int {
	return []int{int(xdr.LedgerEntryTypeAccount)}
}

func (m *AccountTypeMutator) GetStateUpdate(change xdr.LedgerEntryChange) (update StateUpdate) {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		data := change.Created.Data
		switch data.Type {
		case xdr.LedgerEntryTypeAccount:
			update.AccountType = &StateAccountTypeUpdate{
				Address:     data.Account.AccountId.Address(),
				AccountType: data.Account.AccountType,
			}
		}
	case xdr.LedgerEntryChangeTypeUpdated:
		data := change.Updated.Data
		switch data.Type {
		case xdr.LedgerEntryTypeAccount:
			update.AccountType = &StateAccountTypeUpdate{
				Address:     data.Account.AccountId.Address(),
				AccountType: data.Account.AccountType,
			}
		}
	}
	return update
}
