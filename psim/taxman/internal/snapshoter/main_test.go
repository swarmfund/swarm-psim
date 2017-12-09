package snapshoter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestSnapshot(t *testing.T) {
	snapshots := Snapshots{}
	Convey("Given valid snapshoter", t, func() {
		snapshot := New()
		snapshot.Ledger = 123
		// can add
		snapshots.Add(snapshot)
		// can get
		actualSnapshot := snapshots.Get(snapshot.Ledger)
		assert.Equal(t, actualSnapshot, snapshot)
		// add latest
		latestSnapshot := New()
		latestSnapshot.Ledger = snapshot.Ledger + 1
		snapshots.Add(latestSnapshot)

		Convey("Snapshot does not exist", func() {
			snapshot = snapshots.Get(1)
			So(snapshot, ShouldBeNil)
		})
	})
}
