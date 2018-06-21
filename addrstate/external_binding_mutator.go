package addrstate

import (
	"gitlab.com/tokend/go/xdr"
)

func ExternalSystemBindingMutator(systemType int32) func(xdr.LedgerEntryChange) StateUpdate {
	return func(change xdr.LedgerEntryChange) (update StateUpdate) {
		switch change.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			switch change.Created.Data.Type {
			case xdr.LedgerEntryTypeExternalSystemAccountId:
				data := change.Created.Data.ExternalSystemAccountId
				if int32(data.ExternalSystemType) != systemType {
					break
				}
				update.ExternalAccount = &StateExternalAccountUpdate{
					ExternalType: systemType,
					State:        ExternalAccountBindingStateCreated,
					Data:         string(data.Data),
					Address:      data.AccountId.Address(),
				}
			}
		case xdr.LedgerEntryChangeTypeRemoved:
			switch change.Removed.Type {
			case xdr.LedgerEntryTypeExternalSystemAccountId:
				data := change.Removed.ExternalSystemAccountId
				if int32(data.ExternalSystemType) != systemType {
					break
				}
				update.ExternalAccount = &StateExternalAccountUpdate{
					ExternalType: systemType,
					State:        ExternalAccountBindingStateDeleted,
					Address:      data.AccountId.Address(),
				}
			}
		}
		return update
	}
}
