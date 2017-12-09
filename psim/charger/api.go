package charger

import (
	"net/http"

	"gitlab.com/tokend/psim/ape"
	"gitlab.com/tokend/psim/psim/charger/internal/api"
	"gitlab.com/distributed_lab/logan/v3"
)

func (s *Service) Run() {
	r := api.Router(s.log, s.stripe)
	if s.config.Pprof {
		r.Mount("/", ape.DebugRouter())
	}
	s.log.WithFields(logan.F{
		"address": s.listener.Addr().String(),
		"debug":   s.config.Pprof,
	}).Info("listening")
	s.errors <- http.Serve(s.listener, r)
}
