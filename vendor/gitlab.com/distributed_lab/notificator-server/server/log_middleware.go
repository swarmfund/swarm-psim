package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web/mutil"
	"net/http"
)

func LogMiddleware(log *logrus.Entry) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			lw := mutil.WrapWriter(w)

			h.ServeHTTP(lw, r)

			status := lw.Status()
			entry := log.WithField("url", r.URL.String()).WithField("status", lw.Status())
			switch {
			case status >= 100 && status < 400 || status == 429:
				entry.Info()
			case status >= 400 && status < 500:
				entry.Warn()
			case status > 500:
				entry.Error()
			}
		}

		return http.HandlerFunc(fn)
	}
}
