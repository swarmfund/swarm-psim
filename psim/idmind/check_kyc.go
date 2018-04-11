package idmind

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

func (s *Service) processNotChecked(ctx context.Context, request horizon.Request) error {
	txID := getIDMindTXId(request)
	if txID == "" {
		return errors.New("No tx_id in the whole ExternalDetails history, cannot check KYC state without TxID.")
	}
	fields := logan.F{
		"tx_id": txID,
	}

	checkResp, err := s.identityMind.CheckState(txID)
	if err != nil {
		return errors.Wrap(err, "Failed to perform CheckKYCState request by TXId", fields)
	}
	fields["check_response"] = checkResp

	rejectReason, details := s.getCheckRespRejectReason(*checkResp)
	if rejectReason != "" {
		// Need to reject
		blobID, err := s.rejectCheckKYC(ctx, request.ID, request.Hash, *checkResp, rejectReason, details)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest due to reject reason from CheckResponse")
		}

		s.log.WithFields(logan.F{
			"request":            request,
			"reject_blob_id":     blobID,
			"reject_ext_details": details,
		}).Infof("Rejected KYCRequest during Check Task successfully (%s).", rejectReason)
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

		s.emailProcessor.AddEmailAddresses(ctx, s.config.EmailsToNotify)
	}

	return nil
}

func getIDMindTXId(request horizon.Request) string {
	kyc := request.Details.KYC

	var txID string
	for _, extDetails := range kyc.ExternalDetails {
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
