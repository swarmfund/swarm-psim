package server

import (
	"net/http"
)

func ContentTypeMiddleware(contentType string) func(h http.Handler) http.Handler {
	return func (h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
