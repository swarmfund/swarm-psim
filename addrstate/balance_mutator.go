package addrstate

import "gitlab.com/tokend/go/xdr"

func BalanceMutator(asset string) func(xdr.LedgerEntryChange) StateUpdate {
	return func(change xdr.LedgerEntryChange) (update StateUpdate) {
		switch change.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			switch change.Created.Data.Type {
			case xdr.LedgerEntryTypeBalance:
				data := change.Created.Data.Balance
				if string(data.Asset) != asset {
					break
				}
				update.Balance = &StateBalanceUpdate{
					Address: data.AccountId.Address(),
					Balance: data.BalanceId.AsString(),
					Asset:   data.BalanceId.AsString(),
				}
			}
		}
		return update
	}
}
