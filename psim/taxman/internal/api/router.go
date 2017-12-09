package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/api/handlers"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/snapshoter"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

func Router(
	log *logan.Entry, state *state.State, snapshots *snapshoter.Snapshots,
	horizon *horizon.Connector,
) chi.Router {
	r := chi.NewRouter()
	r.Use(LogCtx(log), StateCtx(state), SnapshotsCtx(snapshots))
	r.Post("/", handlers.Verify)
	r.Get("/state", handlers.GetState)
	r.Get("/snapshots", handlers.GetSnapshots)
	r.Get("/health", handlers.Health)
	return r
}

func HorizonCtx(horizon *horizon.Connector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), handlers.HorizonCtxKey, horizon)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func StateCtx(stripe *state.State) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), handlers.StateCtxKey, stripe)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SnapshotsCtx(snapshots *snapshoter.Snapshots) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), handlers.SnapshotsCtxKey, snapshots)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LogCtx(log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), handlers.LogCtxKey, log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
