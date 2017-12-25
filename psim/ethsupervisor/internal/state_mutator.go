package internal

import (
	"fmt"

	"strings"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/addrstate"
)

func StateMutator(baseAsset, depositAsset string) func(change xdr.LedgerEntryChange) addrstate.StateUpdate {

	assetPairUpdate := func(entry *xdr.AssetPairEntry) *int64 {
		if string(entry.Base) != baseAsset || string(entry.Quote) != depositAsset {
			return nil
		}
		price := int64(entry.PhysicalPrice)
		return &price
	}

	return func(change xdr.LedgerEntryChange) addrstate.StateUpdate {
		update := addrstate.StateUpdate{}
		switch change.Type {
		case xdr.LedgerEntryChangeTypeUpdated:
			switch change.Updated.Data.Type {
			case xdr.LedgerEntryTypeAssetPair:
				update.AssetPrice = assetPairUpdate(change.Updated.Data.AssetPair)
			}
		case xdr.LedgerEntryChangeTypeCreated:
			switch change.Created.Data.Type {
			case xdr.LedgerEntryTypeBalance:
				data := change.Created.Data.Balance
				if string(data.Asset) != depositAsset {
					break
				}
				update.Balance = &addrstate.StateBalanceUpdate{
					Address: data.AccountId.Address(),
					Balance: data.BalanceId.AsString(),
				}
			case xdr.LedgerEntryTypeAssetPair:
				update.AssetPrice = assetPairUpdate(change.Created.Data.AssetPair)
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
}
