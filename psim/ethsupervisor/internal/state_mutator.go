package internal

import (
	"fmt"

	"strings"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/addrstate"
)

func StateMutator(change xdr.LedgerEntryChange) addrstate.StateUpdate {
	update := addrstate.StateUpdate{}
	switch change.Type {
	case xdr.LedgerEntryChangeTypeUpdated:
		switch change.Updated.Data.Type {
		case xdr.LedgerEntryTypeAssetPair:
			data := change.Updated.Data.AssetPair
			if data.Base != "ETH" || data.Quote != "SUN" {
				break
			}
			price := int64(data.PhysicalPrice)
			update.AssetPrice = &price
		}
	case xdr.LedgerEntryChangeTypeCreated:
		switch change.Created.Data.Type {
		case xdr.LedgerEntryTypeBalance:
			data := change.Created.Data.Balance
			update.Balance = &addrstate.StateBalanceUpdate{
				Address: data.AccountId.Address(),
				Balance: data.BalanceId.AsString(),
			}
		case xdr.LedgerEntryTypeAssetPair:
			data := change.Created.Data.AssetPair
			if data.Base != "ETH" || data.Quote != "SUN" {
				break
			}
			price := int64(data.PhysicalPrice)
			update.AssetPrice = &price
		case xdr.LedgerEntryTypeExternalSystemAccountId:
			data := change.Created.Data.ExternalSystemAccountId
			switch data.ExternalSystemType {
			case xdr.ExternalSystemTypeEthereum:
				update.Address = &addrstate.StateAddressUpdate{
					Offchain: strings.ToLower(fmt.Sprint("0x", data.Data)),
					Tokend:   data.AccountId.Address(),
				}
			}
		}
	}
	return update
}
