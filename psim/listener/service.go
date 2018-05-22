package listener

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/distributed_lab/running"
	"time"
)

type Service struct {
	config   Config
	listener Listener
	logger   *logan.Entry
}

const (
	defaultRetryTimeIncrement = 1*time.Second
	defaultMaxRetryTime = 30*time.Second
)

func New(config Config, listener Listener, log *logan.Entry) *Service {
	return &Service{
		config:   config,
		listener: listener,
		logger:   log,
	}
}

func (s *Service) Run(ctx context.Context) {
	running.UntilSuccess(ctx, s.logger, conf.ListenerService, s.BroadcastEvents, defaultRetryTimeIncrement, defaultMaxRetryTime)
}
