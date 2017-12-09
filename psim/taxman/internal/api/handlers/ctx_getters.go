package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/snapshoter"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

const (
	StateCtxKey     = "state"
	LogCtxKey       = "log"
	SnapshotsCtxKey = "snapshots"
	HorizonCtxKey   = "horizon"
	SignerCtxKey    = "config"
)

func State(r *http.Request) state.State {
	return *r.Context().Value(StateCtxKey).(*state.State)
}

func Snapshots(r *http.Request) snapshoter.Snapshots {
	return *r.Context().Value(SnapshotsCtxKey).(*snapshoter.Snapshots)
}

func Horizon(r *http.Request) *horizon.Connector {
	return r.Context().Value(HorizonCtxKey).(*horizon.Connector)
}

func Signer(r *http.Request) keypair.KP {
	return r.Context().Value(SignerCtxKey).(keypair.KP)
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(LogCtxKey).(*logan.Entry)
}
