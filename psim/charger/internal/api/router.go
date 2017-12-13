package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stripe/stripe-go/client"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/charger/internal/api/handlers"
)

// TODO Add horizon connector
func Router(log *logan.Entry, stripe *client.API) chi.Router {
	r := chi.NewRouter()
	// TODO Add HorizonCtx middleware
	r.With(StripeCtx(stripe), LogCtx(log))

	r.Post("/stripe", handlers.StripeChargeHandler)
	r.Post("/bitcoin", handlers.BitcoinChargeHandler)

	return r
}

func StripeCtx(stripe *client.API) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.PutStripe(r, stripe)
			next.ServeHTTP(w, r)
		})
	}
}

func LogCtx(log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.PutLog(r, log)
			next.ServeHTTP(w, r)
		})
	}
}

func HorizonCtx(horizon horizon.Connector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO
			next.ServeHTTP(w, r)
		})
	}
}
