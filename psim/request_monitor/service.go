package request_monitor

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

type Service struct {
	config              Config
	logger              *logan.Entry
	connector           *horizon.Connector
	notifier            Notifier
	sentStaleRequestIDs map[uint64]time.Time
}

func New(config Config, notifier Notifier, log *logan.Entry, horizonConnector *horizon.Connector,
) *Service {
	return &Service{
		config:              config,
		logger:              log.WithField("service", conf.ServiceRequestMonitor),
		connector:           horizonConnector,
		notifier:            notifier,
		sentStaleRequestIDs: make(map[uint64]time.Time),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting...")

	running.WithBackOff(
		ctx,
		s.logger,
		conf.ServiceRequestMonitor,
		s.checkRequests,
		s.config.SleepPeriod,
		1*time.Minute,
		1*time.Hour)
}

func (s *Service) checkRequests(ctx context.Context) error {

	ch := s.connector.WithSigner(s.config.Signer).Listener().StreamAllReviewableRequestsOnce(ctx)

	for requestEvent := range ch {
		request, err := requestEvent.Unwrap()
		if err != nil {
			return errors.Wrap(err, "failed to get request")
		}

		if request.State != int32(ReviewableRequestStatePending) {
			delete(s.sentStaleRequestIDs, request.ID)
			continue
		}

		if s.isItTimeToNotify(request.UpdatedAt.Time(), request.Details.RequestType, request.ID) {
			err = s.notifier.Notify(logan.F{
				"request type":  request.Details.RequestTypeName,
				"request state": request.StateName,
				"last update":   request.UpdatedAt.Time().UTC(),
				"request_id":    request.ID,
				"reviewer":      request.Reviewer,
				"requestor":     request.Requestor,
			}, "stale request")

			if err != nil {
				s.logger.WithError(err).Error("failed to notify")
				continue
			}

			s.sentStaleRequestIDs[request.ID] = time.Now().UTC()
		}
	}

	return nil
}

func (s *Service) isUnresolvedBeforeTimeout(updatedAt time.Time, requestType int32) bool {
	timeout := s.getTimeout(requestType)
	elapsedTime := time.Now().Sub(updatedAt)
	return elapsedTime > timeout
}

func (s *Service) isItTimeToNotify(updatedAt time.Time, requestType int32, requestID uint64) bool {
	if lastNotified, ok := s.sentStaleRequestIDs[requestID]; ok {
		notNotified := time.Now().Sub(lastNotified)
		if notNotified < s.config.NotifyPeriod {
			return false
		}
	}

	return s.isUnresolvedBeforeTimeout(updatedAt, requestType)
}

func (s *Service) getTimeout(requestType int32) time.Duration {
	typeString := xdr.ReviewableRequestType(requestType).ShortString()
	timeout, ok := s.config.Requests[typeString]
	if ok {
		return timeout.Timeout
	}
	return s.config.DefaultTimeout

}
