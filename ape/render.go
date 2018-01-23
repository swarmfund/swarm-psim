package ape

import (
	"net/http"

	"strconv"

	"github.com/google/jsonapi"
)

// DEPRECATED Use anything, but not this.
func RenderErr(w http.ResponseWriter, r *http.Request, apiErr *jsonapi.ErrorObject) {
	status, err := strconv.ParseInt(apiErr.Status, 10, 64)
	if err != nil {
		panic(err)
	}

	w.Header().Set("content-type", jsonapi.MediaType)
	w.WriteHeader(int(status))
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{apiErr})
}

// Render will only render slice or pointer.
// DEPRECATED Use anything, but not this.
func Render(w http.ResponseWriter, r *http.Request, res interface{}) {
	w.Header().Set("content-type", jsonapi.MediaType)
	jsonapi.MarshalPayload(w, res)
}
