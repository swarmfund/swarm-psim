package ape

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/pkg/errors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

func Listener(host string, port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
}

func DefaultRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	return r
}

func DebugRouter() chi.Router {
	r := chi.NewRouter()
	r.Handle("/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.Handle("/pprof/*", http.HandlerFunc(pprof.Index))
	return r
}

func InjectPprof(mux chi.Router) {
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	mux.Handle("/debug/pprof/*", http.HandlerFunc(pprof.Index))
}

// Serve wrap `http.Serve` in a channel, so we can select on it
func serve(listener net.Listener, handler http.Handler) (http.Server, chan error) {
	errs := make(chan error)
	server := http.Server{Handler: handler}
	go func() {
		defer close(errs)
		errs <- server.Serve(listener)
	}()
	return server, errs
}

// ListenAndServe will call Shutdown for server once ctx is done.
// NOTE: If Shutdown was successful - this method returns nil.
func ListenAndServe(ctx context.Context, listener net.Listener, handler http.Handler) error {
	server, serveErr := serve(listener, handler)
	select {
	case err := <-serveErr:
		if err != nil && err != http.ErrServerClosed {
			return errors.Wrap(err, "serve failed")
		}
		return nil
	case <-ctx.Done():
		if err := server.Shutdown(context.Background()); err != nil {
			return errors.Wrap(err, "failed to shutdown")
		}
		return nil
	}
}
