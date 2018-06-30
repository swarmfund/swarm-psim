package request_monitor

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"fmt"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdr"
)

type Service struct {
	config Config
	logger *logan.Entry
	connector *horizon.Connector
	ctx context.Context
	stats Stats
}

type RequestTypeToNumber map[xdr.ReviewableRequestType]int

type Stats struct {
	unresolvedRequestIDs []uint64
	requestTypeToNumber RequestTypeToNumber
}

func (st Stats) String() string {
	output := "Number of requests of each type:\n"
	for key, value := range st.requestTypeToNumber {
		output = fmt.Sprintf("%v%v: %v\n", output, key, value)
	}
	output = fmt.Sprintf("%v\nRequests with these IDs haven't been resolved before timeout:\n%v\n", output, st.unresolvedRequestIDs)
	return output
}

func New(config Config, log *logan.Entry, horizonConnector *horizon.Connector) *Service {
	return &Service{
		config: config,
		logger: log.WithField("service", conf.ServiceRequestMonitor),
		connector: horizonConnector,
		stats: Stats{requestTypeToNumber: RequestTypeToNumber{}},
	}
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("starting...")
	s.ctx = ctx

	running.WithBackOff(
		ctx,
		s.logger,
		conf.ServiceRequestMonitor,
		s.worker,
		s.config.SleepPeriod,
		s.config.AbnormalPeriodMin,
		s.config.AbnormalPeriodMax)
}

func (s *Service) worker(ctx context.Context) error {
	s.updateStats()
	s.logger.Info(s.stats)
	return nil
}

func (s *Service) updateStats() {
	ch := s.connector.Listener().StreamAllReviewableRequestsUntilEmptyPage(s.ctx)

	for requestEvent := range ch {
		request, err := requestEvent.Unwrap()
		if err != nil {
			s.logger.Error(errors.Wrap(err, "Failed to unwrap request event"))
			continue
		}

		if s.isUnresolvedBeforeTimeout(request.CreatedAt, request.State) {
			s.stats.unresolvedRequestIDs = append(s.stats.unresolvedRequestIDs, request.ID)
		}

		requestType := xdr.ReviewableRequestType(request.Details.RequestType)
		s.stats.requestTypeToNumber[requestType] += 1
	}
}

func (s *Service) isUnresolvedBeforeTimeout(createdAt time.Time, state int32) bool {
	elapsedTime := time.Now().Sub(createdAt)
	return elapsedTime > s.config.RequestTimeout && state == int32(ReviewableRequestStatePending)
}
