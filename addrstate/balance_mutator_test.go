package addrstate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func TestBalanceMutator(t *testing.T) {
	var accountID xdr.AccountId
	var balanceID xdr.BalanceId
	_ = accountID.SetAddress("GCHIOBEUEOP3WVIZ5ZQIE6HSFGOZFNGB2RQZ2A5TNMF5RWT6M5KD3CFP")
	_ = balanceID.SetString("BAZGUTUOOXXMSW6V7XDW2IM5MIOW4TUYB7BAEJGTLLKCXVWCYRWSSPMP")

	firstLedgerEntry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeBalance,
			Balance: &xdr.BalanceEntry{
				Asset:     "NOTETH",
				AccountId: accountID,
				BalanceId: balanceID,
			},
		},
	}
	firstPayload, err := xdr.MarshalBase64(firstLedgerEntry)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal ledger entry"))
	}

	secondLedgerEntry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeBalance,
			Balance: &xdr.BalanceEntry{
				Asset:     "ETH",
				AccountId: accountID,
				BalanceId: balanceID,
			},
		},
	}
	secondPayload, err := xdr.MarshalBase64(secondLedgerEntry)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal ledger entry"))
	}

	cases := []struct {
		asset    string
		change   regources.LedgerEntryChangeV2
		expected StateUpdate
	}{
		{
			"ETH",
			regources.LedgerEntryChangeV2{
				Effect: int32(xdr.LedgerEntryChangeTypeCreated),
				EntryType: int32(xdr.LedgerEntryTypeBalance),
				Payload: firstPayload,
			},
			StateUpdate{},
		},
		{
			"ETH",
			regources.LedgerEntryChangeV2{
				Effect: int32(xdr.LedgerEntryChangeTypeCreated),
				EntryType: int32(xdr.LedgerEntryTypeBalance),
				Payload: secondPayload,
			},
			StateUpdate{
				Balance: &StateBalanceUpdate{
					Address: accountID.Address(),
					Balance: balanceID.AsString(),
					Asset:   "ETH",
				},
			},
		},
	}
	for i, tc := range cases {
		balanceMutator := BalanceMutator(tc.asset)
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := balanceMutator.GetStateUpdate(tc.change)
			assert.EqualValues(t, tc.expected, got)
		})
	}
}
