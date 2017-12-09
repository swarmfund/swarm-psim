package internal

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/psim/addrstate"
)

func StateMutator(change xdr.LedgerEntryChange) addrstate.StateUpdate {
	return addrstate.StateUpdate{}
}
