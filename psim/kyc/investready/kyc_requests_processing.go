package investready

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/horizon-connector"
)

// TODO timeToSleep to config
// ProcessRequestsInfinitely is blocking method, returns only if ctx is cancelled.
func (s *Service) processRequestsInfinitely(ctx context.Context) {
	for {
		// TODO timeToSleep to config
		timeToSleep := 30 * time.Second

		err := s.processAllRequestsOnce(ctx)
		if err != nil {
			// TODO Add timeToSleep to logs
			s.log.WithError(err).Error("Failed to perform KYCRequests processing iteration. Waiting for the next iteration in a regular mode.")
		} else {
			s.log.Debugf("No more KYC Requests in Horizon, will start from the very beginning, now sleeping for (%s).", timeToSleep.String())
		}

		c := time.After(timeToSleep)
		select {
		case <-ctx.Done():
			return
		case <-c:
			if running.IsCancelled(ctx) {
				return
			}
			continue
		}
	}
}

func (s *Service) processAllRequestsOnce(ctx context.Context) error {
	users, err := s.investReady.ListAllSyncedUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain all SyncedUsers from InvestReady")
	}
	if running.IsCancelled(ctx) {
		return nil
	}
	s.Users = users

	s.kycRequests = s.requestListener.StreamAllKYCRequests(ctx, false)

	running.UntilSuccess(ctx, s.log, "kyc_request_processor", s.listenAndProcessRequest, 5*time.Second, 5*time.Minute)
	return nil
}

func (s *Service) listenAndProcessRequest(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		return true, nil
	case reqEvent, ok := <-s.kycRequests:
		if running.IsCancelled(ctx) {
			return true, nil
		}

		if !ok {
			// No more KYC requests - stopping this iteration.
			return true, nil
		}

		request, err := reqEvent.Unwrap()
		if err != nil {
			return false, errors.Wrap(err, "RequestListener sent error")
		}

		err = s.processRequest(ctx, *request)
		if err != nil {
			return false, errors.Wrap(err, "Failed to process KYC Request", logan.F{
				"request": request,
			})
		}

		return false, nil
	}
}

func (s *Service) processRequest(ctx context.Context, request horizon.Request) error {
	proveErr := proveInterestingRequest(request)
	if proveErr != nil {
		// No need to process the Request for now.

		// I found this log useless
		//s.log.WithField("request", request).WithError(proveErr).Debug("Found not interesting KYC Request.")
		return nil
	}

	// I found this log useless
	s.log.WithField("request", request).Debug("Found interesting KYC Request.")
	//kycReq := request.Details.KYC

	return nil
}
