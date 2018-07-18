package addrstate

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

type BalanceMutator string

func (b *BalanceMutator) GetEffects() []int {
	return []int{int(xdr.LedgerEntryChangeTypeCreated)}
}

func (b *BalanceMutator) GetEntryTypes() []int {
	return []int{int(xdr.LedgerEntryTypeBalance)}
}


func (b *BalanceMutator) GetStateUpdate(change regources.LedgerEntryChangeV2) (update StateUpdate) {
	switch change.EntryType {
	case int32(xdr.LedgerEntryTypeBalance):
		switch change.Effect {
		case int32(xdr.LedgerEntryChangeTypeCreated):
			var ledgerEntry xdr.LedgerEntry
			err := xdr.SafeUnmarshalBase64(change.Payload, &ledgerEntry)
			if err != nil {

			}
			balance := ledgerEntry.Data.MustBalance()
			if string(balance.Asset) != string(*b) {
				break
			}
			update.Balance = &StateBalanceUpdate{
				Address: balance.AccountId.Address(),
				Balance: balance.BalanceId.AsString(),
				Asset:   string(balance.Asset),
			}
		}
	}
	return update
}

