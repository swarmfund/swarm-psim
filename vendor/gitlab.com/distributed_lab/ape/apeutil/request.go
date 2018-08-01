package apeutil

import (
	"bytes"
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

// RequestWithURLParams builds http.Request with body and injects chi urlparams,
// useful for tests
func RequestWithURLParams(body []byte, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	r, _ := http.NewRequest("GET", "/", bytes.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	return r
}
