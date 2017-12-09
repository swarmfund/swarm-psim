package handlers

import (
	"net/http"

	"gitlab.com/tokend/psim/ape"
	"gitlab.com/tokend/psim/psim/taxman/internal/resource"
)

func GetSnapshots(w http.ResponseWriter, r *http.Request) {
	snapshots := resource.Snapshots{}
	for key, snapshot := range Snapshots(r) {
		snapshots = append(snapshots, &resource.Snapshot{
			ID:       key,
			Snapshot: snapshot,
		})
	}
	ape.Render(w, r, snapshots)
}
