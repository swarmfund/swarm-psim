package idmind

import (
	"context"
	"encoding/json"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
)

const (
	RejectorName = "id_mind"
)

func (s *Service) rejectInvalidKYCData(ctx context.Context, requestID uint64, requestHash string, isUSA bool, validationErr error) error {
	var tasksToAdd uint32
	if isUSA {
		tasksToAdd = kyc.TaskUSA
	}

	extDetails := map[string]string{
		"validation_error": validationErr.Error(),
	}

	_, err := s.reject(ctx, requestID, requestHash, nil, s.config.RejectReasons.InvalidKYCData, tasksToAdd, extDetails)
	return err
}

// rejectReason must be absolutely human-readable, we show it to User
func (s *Service) rejectSubmitKYC(ctx context.Context, requestID uint64, requestHash string, idMindResp interface{}, rejectReason string, extDetails map[string]string, isUSA bool) (blobID string, err error) {
	var tasksToAdd uint32
	if isUSA {
		tasksToAdd = kyc.TaskUSA
	}

	return s.reject(ctx, requestID, requestHash, idMindResp, rejectReason, tasksToAdd, extDetails)
}

func (s *Service) rejectCheckKYC(ctx context.Context, requestID uint64, requestHash string, idMindResp interface{}, rejectReason string, extDetails map[string]string) (blobID string, err error) {
	return s.reject(ctx, requestID, requestHash, idMindResp, rejectReason, 0, extDetails)
}

// idMindResp can be nil (in this case blobID in return will be empty)
// extDetails can be nil
func (s *Service) reject(ctx context.Context, requestID uint64, requestHash string, idMindResp interface{}, rejectReason string, tasksToAdd uint32, extDetails map[string]string) (blobID string, err error) {
	if extDetails == nil {
		extDetails = make(map[string]string)
	}

	if idMindResp != nil {
		// Put IDMind response into Blobs.
		idMindRespBB, err := json.Marshal(idMindResp)
		if err != nil {
			return "", errors.Wrap(err, "Failed to marshal provided IDMind response into bytes")
		}

		blobID, err = s.blobSubmitter.SubmitBlob(ctx, "kyc_form", string(idMindRespBB), map[string]string{
			"request_id":   strconv.Itoa(int(requestID)),
			"request_hash": requestHash,
		})
		if err != nil {
			return "", errors.Wrap(err, "Failed to submit Blob via BlobSubmitter")
		}

		extDetails["blob_id"] = blobID
	}

	err = s.requestPerformer.Reject(ctx, requestID, requestHash, tasksToAdd, extDetails, rejectReason, RejectorName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to sign or submit RejectRequest TX")
	}

	return blobID, nil
}
