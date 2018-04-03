package ape

import (
	"net/http"

	"runtime/debug"

	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// RecoverMiddleware by default will just catch handler panic.
// Also it might take use of number of optionally injected arguments like:
// - *logan.Entry to log stacktrace with Error level
// - *jsonapi.ErrorObject to render error body
// - *raven.Client to report exception to Sentry
func RecoverMiddleware(args ...interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					responseRendered := false
					rerr := errors.FromPanic(rvr)
					for _, arg := range args {
						switch v := arg.(type) {
						case *logan.Entry:
							v.
								WithField("stack", string(debug.Stack())).
								WithError(rerr).
								Error("handler panicked")
						case *jsonapi.ErrorObject:
							RenderErr(w, v)
							responseRendered = true
						case *raven.Client:
							packet := raven.NewPacket(
								fmt.Sprint(rvr),
								raven.NewException(
									rerr,
									raven.GetOrNewStacktrace(rerr, 0, 0, nil)),
								raven.NewHttp(r))
							v.Capture(packet, nil)
						}
					}
					if !responseRendered {
						RenderErr(w, problems.InternalError())
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
