package taxman

import (
	"context"

	"gitlab.com/tokend/psim/ape"
	"gitlab.com/tokend/psim/psim/taxman/internal/api"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *Service) API(ctx context.Context) {
	r := api.Router(s.log, s.state, &s.snapshots, s.horizon)
	if s.config.Pprof {
		r.Mount("/", ape.DebugRouter())
	}
	s.log.WithFields(logan.F{
		"address": s.listener.Addr().String(),
		"debug":   s.config.Pprof,
	}).Info("listening")
	if err := ape.ListenAndServe(ctx, s.listener, r); err != nil {
		s.errors <- errors.Wrap(err, "api failed")
		return
	}
	s.log.Debug("api stopped")
}
