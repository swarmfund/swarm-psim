package ape

import (
	"net/http"

	"strconv"

	"github.com/google/jsonapi"
)

func RenderErr(w http.ResponseWriter, r *http.Request, apiErr *jsonapi.ErrorObject) {
	status, err := strconv.ParseInt(apiErr.Status, 10, 64)
	if err != nil {
		panic(err)
	}

	w.Header().Set("content-type", jsonapi.MediaType)
	w.WriteHeader(int(status))
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{apiErr})
}
