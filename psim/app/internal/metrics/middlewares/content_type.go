package middlewares

import "net/http"

//TODO MOVE TO APE
func ContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}
