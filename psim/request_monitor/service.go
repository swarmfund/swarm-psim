package request_monitor

import (
	"context"

	"time"

	"fmt"

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

type RequestTypeToNumber map[xdr.ReviewableRequestType]int

type Stats struct {
	unresolvedRequestIDs []uint64
	requestTypeToNumber  RequestTypeToNumber
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
		10*time.Minute,
		1*time.Hour)
}

func (s *Service) worker(ctx context.Context) error {
	stats := s.generateStats(ctx)
	s.logger.Info(stats)
	return nil
}

func (s *Service) generateStats(ctx context.Context) Stats {
	stats := Stats{requestTypeToNumber: RequestTypeToNumber{}}
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
