package idmind

import (
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

const (
	TxIDExtDetailsKey = "tx_id"
)

// isInterestingRequest checks if service should process request further
func isInterestingRequest(request horizon.Request) bool {
	// request should be pending
	if request.State != kyc.RequestStatePending {
		return false
	}

	// request should be of type UpdateKYC
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeUpdateKyc) {
		return false
	}

	// valid UpdateKYC request has KYC details set
	if request.Details.KYC == nil {
		return false
	}

	// service could process request that already have manual steps resolved
	if request.Details.KYC.PendingTasks&(kyc.TaskFaceValidation|kyc.TaskDocsExpirationDate) != 0 {
		return false
	}

	// service could process only specific tasks
	if request.Details.KYC.PendingTasks&(kyc.TaskSubmitIDMind|kyc.TaskCheckIDMind) == 0 {
		return false
	}

	return true
}
