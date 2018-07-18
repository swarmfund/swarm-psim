package addrstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources"
)

func TestExternalSystemBindingMutator(t *testing.T) {
	var accountID xdr.AccountId
	_ = accountID.SetAddress("GCHIOBEUEOP3WVIZ5ZQIE6HSFGOZFNGB2RQZ2A5TNMF5RWT6M5KD3CFP")

	firstLedgerEntry := &xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeExternalSystemAccountId,
			ExternalSystemAccountId: &xdr.ExternalSystemAccountId{
				ExternalSystemType: 24,
				Data:               "data",
				AccountId:          accountID,
			},
		},
	}
	firstPayload, err := xdr.MarshalBase64(firstLedgerEntry)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal ledger entry"))
	}

	secondLedgerEntry := &xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeExternalSystemAccountId,
			ExternalSystemAccountId: &xdr.ExternalSystemAccountId{
				ExternalSystemType: 42,
				Data:               "data",
				AccountId:          accountID,
			},
		},
	}
	secondPayload, err := xdr.MarshalBase64(secondLedgerEntry)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal ledger entry"))
	}
	firstLedgerKey := &xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeExternalSystemAccountId,
		ExternalSystemAccountId: &xdr.LedgerKeyExternalSystemAccountId{
			AccountId:          accountID,
			ExternalSystemType: 24,
		},
	}

	thridPayload, err := xdr.MarshalBase64(firstLedgerKey)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal ledger key"))
	}
	secondLedgerKey := &xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeExternalSystemAccountId,
		ExternalSystemAccountId: &xdr.LedgerKeyExternalSystemAccountId{
			AccountId:          accountID,
			ExternalSystemType: 42,
		},
	}

	fourthPayload, err := xdr.MarshalBase64(secondLedgerKey)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal ledger key"))
	}

	cases := []struct {
		name       string
		systemType int32
		change     regources.LedgerEntryChangeV2
		expected   StateUpdate
	}{
		{
			"wrong external type create",
			42,
			regources.LedgerEntryChangeV2{
				Effect: int32(xdr.LedgerEntryChangeTypeCreated),
				EntryType: int32(xdr.LedgerEntryTypeExternalSystemAccountId),
				Payload: firstPayload,
			},
			StateUpdate{},
		},
		{
			"valid external type create",
			42,
			regources.LedgerEntryChangeV2{
				Effect: int32(xdr.LedgerEntryChangeTypeCreated),
				EntryType: int32(xdr.LedgerEntryTypeExternalSystemAccountId),
				Payload: secondPayload,
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
			regources.LedgerEntryChangeV2{
				Effect: int32(xdr.LedgerEntryChangeTypeRemoved),
				EntryType: int32(xdr.LedgerEntryTypeExternalSystemAccountId),
				Payload: thridPayload,
			},
			StateUpdate{},
		},
		{
			"valid external type delete",
			42,
			regources.LedgerEntryChangeV2{
				Effect: int32(xdr.LedgerEntryChangeTypeRemoved),
				EntryType: int32(xdr.LedgerEntryTypeExternalSystemAccountId),
				Payload: fourthPayload,
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
		externalSystemBindingMutator := ExternalSystemBindingMutator(tc.systemType)
		t.Run(tc.name, func(t *testing.T) {
			got, err := externalSystemBindingMutator.GetStateUpdate(tc.change)
			if err != nil {
				panic("failed to get statet update in tests")
			}
			assert.EqualValues(t, tc.expected, got)
		})
	}
}
