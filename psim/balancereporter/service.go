package balancereporter

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	horizon "gitlab.com/tokend/horizon-connector"
)

type Service struct {
	config      ServiceConfig
	horizon     *horizon.Connector
	broadcaster *GenericBroadcaster
	logger      *logan.Entry
}

func NewService(config ServiceConfig, horizon *horizon.Connector, broadcaster *GenericBroadcaster, log *logan.Entry) *Service {
	return &Service{
		config:      config,
		horizon:     horizon,
		broadcaster: broadcaster,
		logger:      log,
	}
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Debug("Starting.")
	targets := s.broadcaster.BufferedTargets
	for _, target := range targets {
		defer func() {
			close(target.Data)
		}()
	}
	running.WithBackOff(ctx, s.logger, "event_dispatcher", s.dispatchEvents, 1*time.Hour, 1*time.Second, 30*time.Second)
}

func (s *Service) dispatchEvents(ctx context.Context) error {
	emittedEvents := make(chan BroadcastedReport)
	defer func() {
		close(emittedEvents)
		if r := recover(); r != nil {
			s.logger.WithRecover(r).Error("got panic while closing target channel")
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.WithRecover(r).Error("catched panic while dispatching events")
			}
		}()
		s.broadcaster.BroadcastEvents(ctx, emittedEvents)
	}()

	for _, threshold := range []int64{1000, 10000, 100000, 1000000} {
		if running.IsCancelled(ctx) {
			return nil
		}

		response, err := s.horizon.System().Balances(s.config.AssetCode, threshold)
		if err != nil {
			return errors.Wrap(err, "failed to get balances from horizon")
		}

		asset, err := s.horizon.Assets().ByCode("SWM")
		if err != nil {
			return errors.Wrap(err, "failed to get asset info from horizon")
		}

		if asset == nil {
			return errors.New("SWM asset not found")
		}

		date := time.Now()
		emittedEvents <- BroadcastedReport{response, int64(asset.Issued), threshold, date}
	}

	return nil
}
