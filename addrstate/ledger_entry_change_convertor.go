package addrstate

import (
	"gitlab.com/tokend/regources"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

var ErrUnexpectedEffect = errors.New("unexpected change effect")

func convertLedgerEntryChangeV2(change regources.LedgerEntryChangeV2) (xdr.LedgerEntryChange, error) {
	switch change.Effect {
	case int32(xdr.LedgerEntryChangeTypeRemoved):
		var ledgerKey xdr.LedgerKey
		err := xdr.SafeUnmarshalBase64(change.Payload, &ledgerKey)
		if err != nil {
			return xdr.LedgerEntryChange{}, errors.Wrap(err, "failed to unmarshal ledger key", logan.F{
				"xdr" : change.Payload,
			})
		}
		return xdr.NewLedgerEntryChange(xdr.LedgerEntryChangeType(change.Effect), ledgerKey)
	case int32(xdr.LedgerEntryChangeTypeCreated), int32(xdr.LedgerEntryChangeTypeUpdated):
		var ledgerEntry xdr.LedgerEntry
		err := xdr.SafeUnmarshalBase64(change.Payload, &ledgerEntry)
		if err != nil {
			return xdr.LedgerEntryChange{}, errors.Wrap(err, "failed to unmarshal ledger entry", logan.F{
				"xdr" : change.Payload,
			})
		}
		return xdr.NewLedgerEntryChange(xdr.LedgerEntryChangeType(change.Effect), ledgerEntry)
	default:
		return xdr.LedgerEntryChange{}, errors.Wrap(ErrUnexpectedEffect, "failed to convert ledger entry",
			logan.F{"effect" : change.Effect})
	}
}
