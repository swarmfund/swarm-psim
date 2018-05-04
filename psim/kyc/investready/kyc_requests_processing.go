package investready

import (
	"context"
	"time"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/kyc"
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
	s.users = users

	s.kycRequests = s.requestListener.StreamAllKYCRequests(ctx, false)

	for {
		if running.IsCancelled(ctx) {
			return nil
		}

		noMoreRequests, err := s.listenAndProcessRequest(ctx)
		if err != nil {
			s.log.WithError(err).Error("Failed to process single Request.")
			continue
		}

		if noMoreRequests {
			break
		}
	}
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

	userHash := getInvestReadyUserHash(*request.Details.KYC)
	if userHash == "" {
		// TODO Consider rejecting the request at this point, as we won't be able to Check State any further too.
		// The situation should normally never happen - this service must not be
		// asked to do the Check until the UserHash is put into KYCRequest.
		return errors.New("No user_hash in the whole ExternalDetails history, cannot check AccreditedInvestor status without UserHash.")
	}
	fields := logan.F{
		"user_hash": userHash,
	}

	user := s.findUser(userHash)
	if user == nil {
		return errors.From(errors.New("User with the UserHash from KYCRequest was not found in InvestReady."), fields.Merge(logan.F{
			"known_users": len(s.users),
		}))
	}
	fields["user"] = user

	kycReq := request.Details.KYC
	kycData, err := s.blobDataRetriever.ParseBlobData(*kycReq)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve KYC Blob or parse KYCData")
	}

	err = s.processInvestReadyUser(ctx, request, *user, *kycData)
	if err != nil {
		return errors.Wrap(err, "Failed to process InvestReady User", fields)
	}

	return nil
}

func (s *Service) processInvestReadyUser(ctx context.Context, request horizon.Request, user User, kycData kyc.Data) error {
	logger := s.log.WithFields(logan.F{
		"request":  request,
		"user":     user,
		"kyc_data": kycData,
	})

	checkErr := s.checkUserData(user, kycData)
	if checkErr != "" {
		logger.WithField("check_err", checkErr).Warn("Meet User with mismatched personal info - rejecting KYCRequest.")

		err := s.requestPerformer.Reject(ctx, request.ID, request.Hash, kyc.TaskSuperAdmin, nil, checkErr)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest (because of personal info mismatch)", logan.F{
				"check_err": checkErr,
			})
		}

		return nil
	}

	if user.Status.Message == PendingStatusMessage {
		// Not Accredited yet - pending.
		return nil
	}

	if user.Status.Message == AccreditedStatusMessage {
		err := s.requestPerformer.Approve(ctx, request.ID, request.Hash, 0, kyc.TaskCheckInvestReady, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to approve KYCRequest (InvestReady approved)")
		}

		logger.Info("Approved KYCRequest of approved AccreditedInvestor.")
		return nil
	}

	if user.Status.Message == DeniedStatusMessage {
		err := s.requestPerformer.Reject(ctx, request.ID, request.Hash, 0, nil, "Invest Ready denied.")
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest (InvestReady denied)")
		}

		logger.Info("Rejected KYCRequest (InvestReady denied).")
		return nil
	}

	// TODO Make sure there's no other valid value of Accredited field
	return errors.Errorf("Unexpected Status of the InvestReady User (%s).", user.Status.Message)
}

func (s *Service) findUser(userHash string) *User {
	for _, user := range s.users {
		if user.Hash == userHash {
			return &user
		}
	}

	return nil
}

func (s *Service) checkUserData(user User, kycData kyc.Data) (checkErr string) {
	if user.FirstName != kycData.FirstName {
		return fmt.Sprintf("Expected first name in InvestReady to be (%s), but got (%s)", kycData.FirstName, user.FirstName)
	}

	if user.LastName != kycData.LastName {
		return fmt.Sprintf("Expected last name in InvestReady to be (%s), but got (%s)", kycData.LastName, user.LastName)
	}

	userDob := time.Unix(user.DateOfBirth, 0).UTC()
	kycDataDob := kycData.DateOfBirth.UTC()
	if userDob != kycDataDob {
		return fmt.Sprintf("Expected date of birth in InvestReady to be (%s), but got (%s)",
			kycDataDob.String(), userDob.String())
	}

	// TODO Email?
	return ""
}
