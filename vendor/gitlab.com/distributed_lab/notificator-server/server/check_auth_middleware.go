package server

import (
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/distributed_lab/notificator-server/q"
)

func CheckAuth(r *http.Request) (bool, error) {
	authorization := r.Header.Get("authorization")

	// check just token
	key := strings.TrimPrefix(authorization, "Bearer ")
	pair, err := q.GetQInstance().Auth().ByPublic(key)
	if err != nil {
		return false, err
	}
	return pair != nil, nil
}

func CheckAuthMiddleware(allowUntrusted bool) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allowUntrusted {
				ok, err := CheckAuth(r)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"reason": "internal server error", "msg": "%s"}`, err)
					return
				}

				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintf(w, `{"reason": "signature mismatch"}`)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
