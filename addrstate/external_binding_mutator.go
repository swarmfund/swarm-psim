package addrstate

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

type ExternalSystemBindingMutator int32

func (e *ExternalSystemBindingMutator) GetEffects() []int {
	return []int{int(xdr.LedgerEntryChangeTypeCreated), int(xdr.LedgerEntryChangeTypeRemoved)}
}

func (e *ExternalSystemBindingMutator) GetEntryTypes() []int {
	return []int{int(xdr.LedgerEntryTypeExternalSystemAccountId)}
}

func (e *ExternalSystemBindingMutator) GetStateUpdate(change regources.LedgerEntryChangeV2) (update StateUpdate) {
	switch change.EntryType {
	case int32(xdr.LedgerEntryTypeExternalSystemAccountId):
		switch change.Effect {
		case int32(xdr.LedgerEntryChangeTypeCreated):
			var ledgerEntry xdr.LedgerEntry
			err := xdr.SafeUnmarshalBase64(change.Payload, &ledgerEntry)
			if err != nil {

			}
			data := ledgerEntry.Data.MustExternalSystemAccountId()
			if int32(data.ExternalSystemType) != int32(*e) {
				break
			}
			update.ExternalAccount = &StateExternalAccountUpdate{
				ExternalType: int32(*e),
				State:        ExternalAccountBindingStateCreated,
				Data:         string(data.Data),
				Address:      data.AccountId.Address(),
			}
		case int32(xdr.LedgerEntryChangeTypeRemoved):
			var ledgerKey xdr.LedgerKey
			err := xdr.SafeUnmarshalBase64(change.Payload, &ledgerKey)
			if err != nil {

			}
			data := ledgerKey.MustExternalSystemAccountId()
			if int32(data.ExternalSystemType) != int32(*e) {
				break
			}
			update.ExternalAccount = &StateExternalAccountUpdate{
				ExternalType: int32(*e),
				State:        ExternalAccountBindingStateDeleted,
				Address:      data.AccountId.Address(),
			}
		}
	}
	return update
}
