package addrstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/xdr"
)

func TestExternalSystemBindingMutator(t *testing.T) {
	var accountID xdr.AccountId
	_ = accountID.SetAddress("GCHIOBEUEOP3WVIZ5ZQIE6HSFGOZFNGB2RQZ2A5TNMF5RWT6M5KD3CFP")
	cases := []struct {
		name       string
		systemType int32
		change     xdr.LedgerEntryChange
		expected   StateUpdate
	}{
		{
			"wrong external type create",
			42,
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeCreated,
				Created: &xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeExternalSystemAccountId,
						ExternalSystemAccountId: &xdr.ExternalSystemAccountId{
							ExternalSystemType: 24,
							Data:               "data",
							AccountId:          accountID,
						},
					},
				},
			},
			StateUpdate{},
		},
		{
			"valid external type create",
			42,
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeCreated,
				Created: &xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeExternalSystemAccountId,
						ExternalSystemAccountId: &xdr.ExternalSystemAccountId{
							ExternalSystemType: 42,
							Data:               "data",
							AccountId:          accountID,
						},
					},
				},
			},
			StateUpdate{
				ExternalAccount: &StateExternalAccountUpdate{
					ExternalType: 42,
					State:        ExternalAccountBindingStateCreated,
					Data:         "data",
					Address:      accountID.Address(),
				},
			},
		},
		{
			"wrong external type delete",
			42,
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeRemoved,
				Removed: &xdr.LedgerKey{
					Type: xdr.LedgerEntryTypeExternalSystemAccountId,
					ExternalSystemAccountId: &xdr.LedgerKeyExternalSystemAccountId{
						AccountId:          accountID,
						ExternalSystemType: 24,
					},
				},
			},
			StateUpdate{},
		},
		{
			"valid external type delete",
			42,
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeRemoved,
				Removed: &xdr.LedgerKey{
					Type: xdr.LedgerEntryTypeExternalSystemAccountId,
					ExternalSystemAccountId: &xdr.LedgerKeyExternalSystemAccountId{
						AccountId:          accountID,
						ExternalSystemType: 42,
					},
				},
			},
			StateUpdate{
				ExternalAccount: &StateExternalAccountUpdate{
					ExternalType: 42,
					State:        ExternalAccountBindingStateDeleted,
					Address:      accountID.Address(),
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ExternalSystemBindingMutator(tc.systemType)(tc.change)
			assert.EqualValues(t, tc.expected, got)
		})
	}
}
