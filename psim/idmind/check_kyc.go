package idmind

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

func (s *Service) checkKYCState(ctx context.Context, request horizon.Request) error {
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

	if txID == "" {
		return errors.New("No tx_id in the whole ExternalDetails history, cannot check KYC state without TxID.")
	}
	fields := logan.F{
		"tx_id": txID,
	}

	checkResp, err := s.identityMind.CheckState(txID)
	if err != nil {
		return errors.Wrap(err, "Failed to check state of TX", fields)
	}
	fields["check_response"] = checkResp

	// TODO Maybe additionally determine, whether Application documents were already checked and the result is final (from the 'etr' field).
	switch checkResp.KYCState {
	case UnderReviewKYCState:
		// Not fully reviewed yet, skipping. Will come back to this KYCRequest later.
		return nil
	case AcceptedKYCState:
		err := s.approveCheckKYC(ctx, request.ID, request.Hash)
		if err != nil {
			return errors.Wrap(err, "Failed to approve during Check Task", fields)
		}

		s.log.WithField("request", request).Info("Approved KYCRequest during Check Task successfully.")
		return nil
	case RejectedKYCState:
		err := s.rejectCheckKYC(ctx, request.ID, request.Hash, checkResp, s.config.RejectReasons.KYCStateRejected)
		if err != nil {
			return errors.Wrap(err, "Failed to reject during Check Task", fields)
		}

		s.log.WithField("request", request).Info("Rejected KYCRequest during Check Task successfully.")
		return nil
	}

	return nil
}
