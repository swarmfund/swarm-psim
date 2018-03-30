package idmind

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

// TODO
func (s *Service) checkKYCState(request horizon.Request) error {
	kyc := request.Details.KYC

	var txID string
	for _, extDetails := range kyc.ExternalDetails {
		value, ok := extDetails["tx_id"]
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
	if checkResp.KYCState == UnderReviewKYCState {
		// Not fully reviewed yet, skipping.
		return nil
	}

	// TODO Submit Operations to Core, include state from checkResp

	return errors.New("Not fully implemented.")
}
