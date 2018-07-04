package reporter

import (
	"context"
	"time"

	"gitlab.com/swarmfund/psim/psim/conf"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	horizon "gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

type ServiceConfig struct {
	Signer    keypair.Full `fig:"signer,required"`
	AssetCode string       `fig:"asset_code,required"`
}

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

const (
	defaultServiceRetryTimeIncrement = 1 * time.Second
	defaultMaxServiceRetryTime       = 30 * time.Second
)

func (s *Service) Run(ctx context.Context) {
	running.UntilSuccess(ctx, s.logger, conf.BalanceReporterService, s.dispatchEvents, defaultServiceRetryTimeIncrement, defaultMaxServiceRetryTime)
}

func (s *Service) dispatchEvents(ctx context.Context) (bool, error) {
	emittedEvents := make(chan BroadcastedReport)
	defer func() {
		close(emittedEvents)
		if r := recover(); r != nil {
			s.logger.Error("got panic while closing target channel")
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("catched panic while dispatching events")
			}
		}()
		s.broadcaster.BroadcastEvents(ctx, emittedEvents)
	}()

	ticks := time.Tick(1 * time.Second)
ticklabel:
	for _ = range ticks {
		for _, tx := range []int64{1000, 10000, 100000, 1000000} {
			if running.IsCancelled(ctx) {
				break ticklabel
			}

			response, err := s.horizon.System().Balances(s.config.AssetCode, tx)
			if err != nil {
				s.logger.WithError(err).Error("failed to get balances from horizon")
				continue ticklabel
			}
			date := time.Now()
			assets, err := s.horizon.Assets().ByCode("SWM")
			if err != nil {
				s.logger.WithError(err).Error("failed to get asset info from horizon")
				continue ticklabel
			}
			if assets == nil {
				s.logger.Error("asset not found")
			}
			emittedEvents <- BroadcastedReport{response, int64(assets.Issued), tx, &date}
		}
	}

	s.logger.Debug("ticker cycle is stopped")
	return false, nil
}
