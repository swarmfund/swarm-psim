package resource

import (
	"gitlab.com/swarmfund/psim/psim/taxman/internal/snapshoter"
)

type Snapshots []*Snapshot

type Snapshot struct {
	ID       int64               `jsonapi:"primary,snapshoter"`
	Snapshot snapshoter.Snapshot `jsonapi:"attr,snapshoter"`
}
