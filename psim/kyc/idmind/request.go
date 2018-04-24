package idmind

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

const (
	RequestStatePending int32 = 1

	TxIDExtDetailsKey = "tx_id"
)

// ProveInterestingRequest returns non-nil error if the provided Request
// doesn't need to be considered by this Service.
func proveInterestingRequest(request horizon.Request) error {
	if request.State != RequestStatePending {
		// State is not pending
		return errors.Errorf("Invalid Request State (%d) expected Pending(%d).", request.State, RequestStatePending)
	}

	details := request.Details

	if details.RequestType != int32(xdr.ReviewableRequestTypeUpdateKyc) {
		return errors.Errorf("Invalid RequestType (%d) expected KYC(%d).", details.RequestType, xdr.ReviewableRequestTypeUpdateKyc)
	}

	kycRequest := details.KYC

	if kycRequest == nil {
		return errors.New("KYC struct in the ReviewableRequest is nil.")
	}

	return proveInterestingMask(kycRequest.PendingTasks)
}

func proveInterestingMask(pendingTasks uint32) error {
	if pendingTasks&(kyc.TaskFaceValidation|kyc.TaskDocsExpirationDate) != 0 {
		// Some manual check hasn't completed - too early to process this request.
		return errors.New("Either FaceValidation or DocsExpirationDate hasn't been checked yet - too early to process this Request.")
	}

	if pendingTasks&(kyc.TaskSubmitIDMind|kyc.TaskCheckIDMind) == 0 {
		return errors.New("No pending tasks for me - ignoring this Request.")
	}

	return nil
}
