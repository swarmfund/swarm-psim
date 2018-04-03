package idmind

import (
	"fmt"

	"context"

	"encoding/json"

	"strconv"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

const (
	RequestStatePending int32 = 1

	TaskSuperAdmin         uint32 = 1
	TaskFaceValidation     uint32 = 2
	TaskDocsExpirationDate uint32 = 4
	TaskSubmitIDMind       uint32 = 8
	TaskCheckIDMind        uint32 = 16
	TaskUSA                uint32 = 32

	TxIDExtDetailsKey = "tx_id"
)

func proveInterestingRequest(request horizon.Request) error {
	if request.State != RequestStatePending {
		// State is not pending
		return errors.Errorf("Invalid Request State (%d) expected Pending(%d).", request.State, RequestStatePending)
	}

	details := request.Details

	if details.RequestType != int32(xdr.ReviewableRequestTypeUpdateKyc) {
		return errors.Errorf("Invalid RequestType (%d) expected KYC(%d).", details.RequestType, xdr.ReviewableRequestTypeUpdateKyc)
	}

	kyc := details.KYC

	if kyc == nil {
		return errors.New("KYC struct in the ReviewableRequest is nil.")
	}

	return proveInterestingMask(kyc.PendingTasks)
}

func proveInterestingMask(pendingTasks uint32) error {
	if pendingTasks&(TaskFaceValidation|TaskDocsExpirationDate) != 0 {
		// Some manual check hasn't completed - too early to process this request.
		return errors.New("Either FaceValidation or DocsExpirationDate hasn't been checked yet - too early to process this Request.")
	}

	if pendingTasks&(TaskSubmitIDMind|TaskCheckIDMind) == 0 {
		return errors.New("No pending tasks for me - ignoring this Request.")
	}

	return nil
}

func (s *Service) rejectInvalidKYCData(ctx context.Context, requestID uint64, requestHash string, isUSA bool, validationErr error) error {
	var tasksToAdd uint32
	if isUSA {
		tasksToAdd = TaskUSA
	}

	extDetails := map[string]string{
		"validation_error": validationErr.Error(),
	}

	return s.reject(ctx, requestID, requestHash, nil, s.config.RejectReasons.InvalidKYCData, tasksToAdd, extDetails)
}

// rejectReason must be absolutely human-readable, we show it to User
func (s *Service) rejectSubmitKYC(ctx context.Context, requestID uint64, requestHash string, idMindResp interface{}, rejectReason string, isUSA bool) error {
	var tasksToAdd uint32
	if isUSA {
		tasksToAdd = TaskUSA
	}

	return s.reject(ctx, requestID, requestHash, idMindResp, rejectReason, tasksToAdd, nil)
}

func (s *Service) rejectCheckKYC(ctx context.Context, requestID uint64, requestHash string, idMindResp interface{}, rejectReason string) error {
	return s.reject(ctx, requestID, requestHash, idMindResp, rejectReason, 0, nil)
}

// idMindResp can be nil
// extDetails can be nil
func (s *Service) reject(ctx context.Context, requestID uint64, requestHash string, idMindResp interface{}, rejectReason string, tasksToAdd uint32, extDetails map[string]string) error {
	if extDetails == nil {
		extDetails = make(map[string]string)
	}

	if idMindResp != nil {
		// Pu IDMind response into Blob.
		idMindRespBB, err := json.Marshal(idMindResp)
		if err != nil {
			return errors.Wrap(err, "Failed to marshal provided IDMind response into bytes")
		}

		blobID, err := s.blobsConnector.SubmitBlob(ctx, "kyc_form", string(idMindRespBB), map[string]string{
			"request_id":   strconv.Itoa(int(requestID)),
			"request_hash": requestHash,
		})
		if err != nil {
			return errors.Wrap(err, "Failed to submit Blob via BlobsConnector")
		}

		extDetails["blob_id"] = blobID
	}

	extDetailsBB, err := json.Marshal(extDetails)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal externalDetails", logan.F{"ext_details": extDetails})
	}

	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.Source).Op(xdrbuild.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionReject,
		Details: xdrbuild.UpdateKYCDetails{
			TasksToAdd:      tasksToAdd,
			TasksToRemove:   0,
			ExternalDetails: string(extDetailsBB),
		},
		Reason: rejectReason,
	}).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal signed Envelope")
	}

	submitResult := s.txSubmitter.Submit(ctx, signedEnvelope)
	if submitResult.Err != nil {
		return errors.Wrap(submitResult.Err, "Error submitting signed Envelope to Horizon", logan.F{
			"submit_result": submitResult,
		})
	}

	return nil
}

func (s *Service) approveSubmitKYC(ctx context.Context, requestID uint64, requestHash, txID string, isUSA bool) error {
	var tasksToAdd = TaskCheckIDMind
	if isUSA {
		tasksToAdd = tasksToAdd | TaskUSA
	}

	extDetails := fmt.Sprintf(`{"%s":"%s"}`, TxIDExtDetailsKey, txID)
	return s.approve(ctx, requestID, requestHash, tasksToAdd, TaskSubmitIDMind, extDetails)
}

func (s *Service) approveCheckKYC(ctx context.Context, requestID uint64, requestHash string) error {
	// TODO In future we will probably need to add some Task at this point (e.g. for some particular admin to make some final review of final accept, or whatever)
	return s.approve(ctx, requestID, requestHash, 0, TaskCheckIDMind, "{}")
}

func (s *Service) approve(ctx context.Context, requestID uint64, requestHash string, tasksToAdd, tasksToRemove uint32, extDetails string) error {
	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.Source).Op(xdrbuild.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.UpdateKYCDetails{
			TasksToAdd:      tasksToAdd,
			TasksToRemove:   tasksToRemove,
			ExternalDetails: extDetails,
		},
		Reason: "",
	}).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal signed Envelope")
	}

	submitResult := s.txSubmitter.Submit(ctx, signedEnvelope)
	if submitResult.Err != nil {
		return errors.Wrap(submitResult.Err, "Error submitting signed Envelope to Horizon", logan.F{
			"submit_result": submitResult,
		})
	}

	return nil
}
