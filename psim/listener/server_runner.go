package listener

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

// RunServer is a blocking method which creates http.Server and runs it until ctx is cancelled.
func RunServer(ctx context.Context, log *logan.Entry, handler http.Handler, config Config) {
	var server *http.Server
	go running.UntilSuccess(ctx, log, "listening_server", func(ctx context.Context) (bool, error) {
		server = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
			Handler:      handler,
			WriteTimeout: config.Timeout,
		}

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return false, errors.Wrap(err, "Failed to ListenAndServe (Server stopped with error)")
		}

		return false, nil
	}, time.Second, time.Hour)

	<-ctx.Done()

	shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(shutdownCtx)
	log.Info("Server stopped cleanly.")
}
