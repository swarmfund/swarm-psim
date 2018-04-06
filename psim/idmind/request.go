package idmind

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
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
	TaskNonLatinDoc        uint32 = 64

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
