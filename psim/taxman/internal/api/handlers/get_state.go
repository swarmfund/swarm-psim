package handlers

import (
	"net/http"

	"gitlab.com/tokend/psim/ape"
	"gitlab.com/tokend/psim/psim/taxman/internal/resource"
)

func GetState(w http.ResponseWriter, r *http.Request) {
	state := resource.StateResource{
		State(r),
	}
	ape.Render(w, r, &state)
}
