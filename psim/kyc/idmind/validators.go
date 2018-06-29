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

	// service should process request that already have manual steps resolved
	if request.Details.KYC.PendingTasks&(kyc.TaskFaceValidation|kyc.TaskDocsExpirationDate) != 0 {
		return false
	}

	// service can process only specific tasks
	if request.Details.KYC.PendingTasks&(kyc.TaskSubmitIDMind|kyc.TaskCheckIDMind) == 0 {
		return false
	}

	// valid UpdateKYC request has KYC details set
	if request.Details.KYC == nil {
		// TODO consider rejecting invalid requests
		return false
	}

	// we expected specific KYC data format to be able perform validation
	if request.Details.KYC.KYCDataStruct.BlobID == "" {
		// TODO consider rejecting invalid requests
		return false
	}

	return true
}

func isBlobValid(blob *horizon.Blob) bool {
	// blob should exists
	if blob == nil {
		return false
	}

	// service expects specific blob type
	if blob.Type != kyc.KYCFormBlobType {
		return false
	}

	return true
}

func isGeneral(account *horizon.Account) bool {
	return account.AccountTypeI == int32(xdr.AccountTypeGeneral)
}

func isNotVerified(account *horizon.Account) bool {
	return account.AccountTypeI == int32(xdr.AccountTypeNotVerified)
}

func isUpdateToGeneral(request horizon.Request) bool {
	return request.Details.KYC.AccountTypeToSet.Int == int(xdr.AccountTypeGeneral)
}

func isUpdateToVerified(request horizon.Request) bool {
	panic("not implemented")
}
