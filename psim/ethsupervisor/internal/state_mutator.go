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
			case xdr.ExternalSystemTypeEthereum:
				update.Address = &addrstate.StateAddressUpdate{
					Offchain: "0xd96c70a7DC0BBEdB9dAb293a7b6a3557B073394e",
					Tokend:   data.AccountId.Address(),
				}
			}
		}
	}
	return update
}
