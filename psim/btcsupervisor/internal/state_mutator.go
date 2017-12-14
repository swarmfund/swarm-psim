package internal

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/addrstate"
)

func StateMutator(change xdr.LedgerEntryChange) addrstate.StateUpdate {
	update := addrstate.StateUpdate{}

	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		switch change.Created.Data.Type {
		case xdr.LedgerEntryTypeExternalSystemAccountId:
			data := change.Created.Data.ExternalSystemAccountId

			switch data.ExternalSystemType {
			case xdr.ExternalSystemTypeBitcoin:
				update.Address = &addrstate.StateAddressUpdate{
					Offchain: string(data.Data),
					Tokend:   data.AccountId.Address(),
				}
			}
		}
	}
	return update
}
