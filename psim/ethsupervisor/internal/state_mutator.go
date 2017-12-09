package internal

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/addrstate"
)

func StateMutator(change xdr.LedgerEntryChange) addrstate.StateUpdate {
	return addrstate.StateUpdate{}
}
