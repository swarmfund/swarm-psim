package addrstate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/xdr"
)

func TestBalanceMutator(t *testing.T) {
	var accountID xdr.AccountId
	var balanceID xdr.BalanceId
	_ = accountID.SetAddress("GCHIOBEUEOP3WVIZ5ZQIE6HSFGOZFNGB2RQZ2A5TNMF5RWT6M5KD3CFP")
	_ = balanceID.SetString("BAZGUTUOOXXMSW6V7XDW2IM5MIOW4TUYB7BAEJGTLLKCXVWCYRWSSPMP")

	cases := []struct {
		asset    string
		change   xdr.LedgerEntryChange
		expected StateUpdate
	}{
		{
			"ETH",
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeCreated,
				Created: &xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeBalance,
						Balance: &xdr.BalanceEntry{
							Asset:     "NOTETH",
							AccountId: accountID,
							BalanceId: balanceID,
						},
					},
				},
			},
			StateUpdate{},
		},
		{
			"ETH",
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeCreated,
				Created: &xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeBalance,
						Balance: &xdr.BalanceEntry{
							Asset:     "ETH",
							AccountId: accountID,
							BalanceId: balanceID,
						},
					},
				},
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
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := BalanceMutator(tc.asset)(tc.change)
			assert.EqualValues(t, tc.expected, got)
		})
	}
}
