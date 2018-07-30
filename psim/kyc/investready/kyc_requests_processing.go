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

const (
	RejectorName       = "invest_ready"
	DeniedRejectReason = "Invest Ready denied."
)

func (s *Service) processAllRequestsOnce(ctx context.Context) error {
	users, err := s.investReady.ListAllSyncedUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain all SyncedUsers from InvestReady")
	}
	if running.IsCancelled(ctx) {
		return nil
	}
	s.syncedUserHashes = users

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

	s.log.Debug("No more KYC Requests in Horizon, will start from the very beginning, now sleeping.")
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

	user := s.findInvestReadyUser(userHash)
	if user == nil {
		return errors.From(errors.New("User with the UserHash from KYCRequest was not found in InvestReady."), fields.Merge(logan.F{
			"known_users": len(s.syncedUserHashes),
		}))
	}
	fields["user"] = user

	kycData, err := s.getBlobKYCData(request)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve KYCData of the Request (from Blob)")
	}

	err = s.processInvestReadyUser(ctx, request, *user, *kycData)
	if err != nil {
		return errors.Wrap(err, "Failed to process InvestReady User", fields)
	}

	return nil
}

func (s *Service) getBlobKYCData(request horizon.Request) (*kyc.Data, error) {
	kycReq := request.Details.KYC
	if kycReq == nil {
		return nil, errors.New("KYCRequest in the Request is nil.")
	}
	fields := logan.F{
		"blob_id": kycReq.KYCDataStruct.BlobID,
	}

	blob, err := s.blobsConnector.Blob(kycReq.KYCDataStruct.BlobID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Blob", fields)
	}
	fields["blob"] = blob

	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse KYCData from the Blob")
	}

	return kycData, nil
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

		err := s.requestPerformer.Reject(ctx, request.ID, request.Hash, kyc.TaskSuperAdmin, nil, checkErr, RejectorName)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest (because of personal info mismatch)", logan.F{
				"check_err": checkErr,
			})
		}

		return nil
	}

	switch user.Status.Message {
	case PendingStatusMessage, NoPendingVerificationsStatusMessage:
		// Not Accredited yet - pending.
		return nil
	case AccreditedStatusMessage:
		err := s.requestPerformer.Approve(ctx, request.ID, request.Hash, 0, kyc.TaskCheckInvestReady, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to approve KYCRequest (InvestReady approved)")
		}

		logger.Info("Approved KYCRequest of approved AccreditedInvestor.")
		return nil
	case DeniedStatusMessage:
		err := s.requestPerformer.Reject(ctx, request.ID, request.Hash, 0, nil, DeniedRejectReason, RejectorName)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest (InvestReady denied)")
		}

		logger.Info("Rejected KYCRequest (InvestReady denied).")
		return nil
	default:
		// TODO Make sure there's no other valid value of Accredited field
		return errors.Errorf("Unexpected Status of the InvestReady User (%s).", user.Status.Message)
	}
}

func (s *Service) findInvestReadyUser(userHash string) *User {
	for _, user := range s.syncedUserHashes {
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
