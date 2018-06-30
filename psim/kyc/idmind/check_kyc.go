package idmind

import (
	"context"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/horizon-connector"
)

// TODO Try to refactor - make method shorter.
func (s *Service) processNotChecked(ctx context.Context, request horizon.Request) error {
	txID := getIDMindTXId(request)
	if txID == "" {
		// TODO Consider rejecting the request at this point, as we won't be able to Check State any further too.
		// The situation should normally never happen - this service must not be asked to do the Check until the Submit is done.
		return errors.New("No tx_id in the whole ExternalDetails history, cannot check KYC state without TxID.")
	}
	fields := logan.F{
		"tx_id": txID,
	}

	logger := s.log.WithFields(logan.F{"request": request})

	checkResp, err := s.identityMind.CheckState(txID)
	if err != nil {
		if err == ErrAppNotFound {
			err = s.requestPerformer.Approve(ctx, request.ID, request.Hash, kyc.TaskSuperAdmin, kyc.TaskCheckIDMind, map[string]string{
				"id_mind_check_result": fmt.Sprintf("Applications with the TxID (%s) was not found in IDMind.", txID),
			})
			if err != nil {
				return errors.Wrap(err, "Failed to approve KYCRequest during Check State (due to AppNotFound IDMind error)", fields)
			}

			logger.Infof("Approved (due to AppNotFound IDMind error) KYCRequest during Check Task successfully.")
			return nil
		} else {

		}
		return errors.Wrap(err, "Failed to perform CheckKYCState request by TXId to IDMind", fields)
	}
	fields["check_response"] = checkResp

	rejectReason, details := s.getCheckRespRejectReason(*checkResp)
	if rejectReason != "" {
		blobID, err := s.reject(ctx, request, *checkResp, rejectReason, details)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest due to reject reason from CheckResponse")
		}

		logger.WithFields(logan.F{
			"reject_blob_id":     blobID,
			"reject_reason":      rejectReason,
			"reject_ext_details": details,
		}).Info("Rejected KYCRequest during Check Task successfully.")
		return nil
	}

	if checkResp.KYCState == AcceptedKYCState {
		err := s.approveCheckKYC(ctx, request.ID, request.Hash)
		if err != nil {
			return errors.Wrap(err, "Failed to approve during Check Task", fields)
		}

		s.log.WithField("request", request).Info("Approved KYCRequest during Check Task successfully.")
		return nil
	}

	// Not fully reviewed yet
	if checkResp.IsManualReview() {
		s.log.WithField("request", request).WithFields(fields).
			Info("Result of immediate response for Application submit is ManualReview, adding notification emails to sending.")

		s.adminNotifyEmails.AddEmailAddresses(ctx, s.config.AdminNotifyEmailsConfig.Subject,
			s.config.AdminNotifyEmailsConfig.Message, s.config.AdminEmailsToNotify)
	}

	return nil
}

func getIDMindTXId(request horizon.Request) string {
	kycReq := request.Details.KYC

	var txID string
	for _, extDetails := range kycReq.ExternalDetails {
		value, ok := extDetails[TxIDExtDetailsKey]
		if !ok {
			// No 'tx_id' key in these externalDetails.
			continue
		}

		txID, ok = value.(string)
		if !ok {
			// Must never happen, but just in case.
			// Maybe we need to log this shit here, if it happens..
			continue
		}
	}

	return txID
}

func (s *Service) getCheckRespRejectReason(checkAppResponse CheckApplicationResponse) (rejectReason string, details map[string]string) {
	if checkAppResponse.KYCState != RejectedKYCState {
		// Not rejected
		return "", nil
	}

	firedRules := checkAppResponse.EDNAScoreCard.FraudPolicyEvaluation.FiredRules
	if len(firedRules) > 0 {
		details = make(map[string]string)
		for _, rule := range firedRules {
			details[rule.Name] = rule.Description
		}
	}

	return s.config.RejectReasons.KYCStateRejected, details
}
