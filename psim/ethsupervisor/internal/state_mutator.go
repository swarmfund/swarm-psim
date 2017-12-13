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
					// TODO unhardcode
					// TODO ensure 0x
					Offchain: "0xBDf4fdB5B70B65791C0A97527796f607Cf846f18",
					Tokend:   data.AccountId.Address(),
				}
			}
		}
	}
	return update
}
