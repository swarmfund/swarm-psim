package internal

import (
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/tokend/go/xdr"
)

const (
	ExternalSystemTypeBitcoin = 1
)

func StateMutator(change xdr.LedgerEntryChange) addrstate.StateUpdate {
	update := addrstate.StateUpdate{}
	switch change.Type {
	case xdr.LedgerEntryChangeTypeUpdated:
		switch change.Updated.Data.Type {
		case xdr.LedgerEntryTypeAssetPair:
			data := change.Updated.Data.AssetPair
			if data.Base != "BTC" || data.Quote != "SUN" {
				break
			}
			price := int64(data.PhysicalPrice)
			update.AssetPrice = &price
		}
	case xdr.LedgerEntryChangeTypeCreated:
		switch change.Created.Data.Type {
		case xdr.LedgerEntryTypeAssetPair:
			data := change.Created.Data.AssetPair
			if data.Base != "BTC" || data.Quote != "SUN" {
				break
			}
			price := int64(data.PhysicalPrice)
			update.AssetPrice = &price
		case xdr.LedgerEntryTypeExternalSystemAccountId:
			data := change.Created.Data.ExternalSystemAccountId

			switch data.ExternalSystemType {
			case ExternalSystemTypeBitcoin:
				update.Address = &addrstate.StateAddressUpdate{
					Offchain: string(data.Data),
					Tokend:   data.AccountId.Address(),
				}
			}
		}
	}
	return update
}
