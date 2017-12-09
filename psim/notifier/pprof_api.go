package notifier

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/psim/ape"
)

func (s *Service) servePProfAPI(ctx context.Context) {
	if !s.Pprof {
		return
	}

	host := "localhost"
	if s.Host == "" {
		host = s.Host
	}

	listener, err := ape.Listener(host, s.Port)
	if err != nil {
		s.logger.WithError(err).Warn("Cant init listener. PProf api disabled")
		return
	}

	router := ape.DebugRouter()
	s.logger.WithField("address", listener.Addr().String()).Info("PProf api listen")

	if err := ape.ListenAndServe(ctx, listener, router); err != nil {
		s.errors <- errors.Wrap(err, "api failed")
		return
	}

	s.logger.Debug("PProf api stopped")
}
