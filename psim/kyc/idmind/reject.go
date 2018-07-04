package idmind

import (
	"context"
	"encoding/json"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
)

const (
	RejectorName = "id_mind"
)

func (s *Service) rejectInvalidKYCData(ctx context.Context, request horizon.Request, validationErr error) error {
	extDetails := map[string]string{
		"validation_error": validationErr.Error(),
	}

	_, err := s.reject(ctx, request, nil, s.config.RejectReasons.InvalidKYCData, extDetails)
	return err
}

func (s *Service) rejectRequest(ctx context.Context, request horizon.Request, rejectReason string, externalDetails map[string]string) error {
	return s.requestPerformer.Reject(ctx, request.ID, request.Hash, 0, externalDetails, rejectReason, RejectorName)
}

// idMindResp can be nil (in this case blobID in return will be empty)
// extDetails can be nil
func (s *Service) reject(ctx context.Context, request horizon.Request, idMindResp interface{}, rejectReason string, extDetails map[string]string) (blobID string, err error) {
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
			"request_id":   strconv.Itoa(int(request.ID)),
			"request_hash": request.Hash,
		})
		if err != nil {
			return "", errors.Wrap(err, "Failed to submit Blob via BlobSubmitter")
		}

		extDetails["blob_id"] = blobID
	}

	err = s.requestPerformer.Reject(ctx, request.ID, request.Hash, 0, extDetails, rejectReason, RejectorName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to sign or submit RejectRequest TX")
	}

	return blobID, nil
}
