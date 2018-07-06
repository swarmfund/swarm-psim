package request_monitor

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

type Service struct {
	config    Config
	logger    *logan.Entry
	connector *horizon.Connector
}

type Stats struct {
	unresolvedRequestIDs []uint64
	requestTypeToNumber  map[xdr.ReviewableRequestType]int
}

func New(config Config, log *logan.Entry, horizonConnector *horizon.Connector) *Service {
	return &Service{
		config:    config,
		logger:    log.WithField("service", conf.ServiceRequestMonitor),
		connector: horizonConnector,
	}
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting...")

	running.WithBackOff(
		ctx,
		s.logger,
		conf.ServiceRequestMonitor,
		s.worker,
		s.config.SleepPeriod,
		1*time.Minute,
		1*time.Hour)
}

func (s *Service) worker(ctx context.Context) error {
	stats := s.generateStats(ctx)
	s.logger.WithFields(logan.F{
		"number_of_requests_by_type": stats.requestTypeToNumber,
		"unresolved_requests_IDs":    stats.unresolvedRequestIDs,
	}).Info("Statistics")
	return nil
}

func (s *Service) generateStats(ctx context.Context) Stats {
	stats := makeEmptyStats()
	ch := s.connector.Listener().StreamAllReviewableRequestsOnce(ctx)

	for requestEvent := range ch {
		request, err := requestEvent.Unwrap()
		if err != nil {
			s.logger.WithError(err)
			continue
		}

		if s.isUnresolvedBeforeTimeout(request.CreatedAt, request.State) {
			stats.unresolvedRequestIDs = append(stats.unresolvedRequestIDs, request.ID)
		}

		requestType := xdr.ReviewableRequestType(request.Details.RequestType)
		stats.requestTypeToNumber[requestType] += 1
	}

	return stats
}

func (s *Service) isUnresolvedBeforeTimeout(createdAt time.Time, state int32) bool {
	elapsedTime := time.Now().Sub(createdAt)
	return elapsedTime > s.config.RequestTimeout && state == int32(ReviewableRequestStatePending)
}

func makeEmptyStats() Stats {
	stats := Stats{requestTypeToNumber: make(map[xdr.ReviewableRequestType]int)}
	for _, requestType := range xdr.ReviewableRequestTypeAll {
		stats.requestTypeToNumber[requestType] = 0
	}
	return stats
}
