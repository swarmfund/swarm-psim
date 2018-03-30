package idmind

import (
	"fmt"

	"context"

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

// TODO Blob submit
// rejectReason must be absolutely human-readable, we show it to User
func (s *Service) reject(ctx context.Context, requestID uint64, requestHash string, idMindResp interface{}, rejectReason string) error {
	// TODO Submit Blob with idMindRespBB to API, get blobID
	//idMindRespBB, err := json.Marshal(idMindResp)
	//if err != nil {
	//	return errors.Wrap(err, "Failed to marshal provided IDMind response into bytes")
	//}

	// TODO Submit Blob with idMindRespBB to API, get blobID

	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.Source).Op(xdrbuild.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionReject,
		Details: xdrbuild.UpdateKYCDetails{
			TasksToAdd:    0,
			TasksToRemove: 0,
			// FIXME
			//ExternalDetails: fmt.Sprintf(`{"blob_id":"%s"}`, blobID),
			ExternalDetails: fmt.Sprintf(`{"blob_id":"%s"}`, "blob_id_stub"),
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

func (s *Service) approveSubmitKYC(ctx context.Context, requestID uint64, requestHash, txID string) error {
	extDetails := fmt.Sprintf(`{"%s":"%s"}`, TxIDExtDetailsKey, txID)
	return s.approve(ctx, requestID, requestHash, TaskCheckIDMind, TaskSubmitIDMind, extDetails)
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
