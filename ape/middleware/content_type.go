package middleware

import "net/http"

func RestrictToContentType(contentTypes ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("content-type")
			for _, allowed := range contentTypes {
				if contentType == allowed {
					h.ServeHTTP(w, r)
				}
			}
			http.Error(w, "", http.StatusUnsupportedMediaType)
		}
		return http.HandlerFunc(fn)
	}
}
