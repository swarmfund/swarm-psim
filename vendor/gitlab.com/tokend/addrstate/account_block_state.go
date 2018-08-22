package addrstate

import "gitlab.com/tokend/go/xdr"

type AccountBlockReasonsMutator struct{}

func (m AccountBlockReasonsMutator) GetEffects() []int {
	return []int{int(xdr.LedgerEntryChangeTypeUpdated)}
}

func (m AccountBlockReasonsMutator) GetEntryTypes() []int {
	return []int{int(xdr.LedgerEntryTypeAccount)}
}

func (m *AccountBlockReasonsMutator) GetStateUpdate(change xdr.LedgerEntryChange) (update StateUpdate) {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeUpdated:
		data := change.Updated.Data
		switch data.Type {
		case xdr.LedgerEntryTypeAccount:
			update.AccountBlockReasons = &StateAccountBlockReasonsUpdate{
				Address:      data.Account.AccountId.Address(),
				BlockReasons: uint32(data.Account.BlockReasons),
			}
		}
	}
	return update
}
