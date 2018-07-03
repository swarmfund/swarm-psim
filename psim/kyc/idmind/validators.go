package idmind

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

const (
	TxIDExtDetailsKey = "tx_id"
)

var (
	errMissingDetails  = errors.New("KYC details are missing")
	errInvalidDetails  = errors.New("KYC details are invalid")
	errBlobNotFound    = errors.New("blob does not exist")
	errInvalidBlobType = errors.New("invalid blob type")
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

	return true
}

func isValidRequest(request horizon.Request) error {
	// valid UpdateKYC request has KYC details set
	if request.Details.KYC == nil {
		return errMissingDetails
	}

	// we expected specific KYC data format to be able to perform validation
	if request.Details.KYC.KYCDataStruct.BlobID == "" {
		return errInvalidDetails
	}

	return nil
}

func isBlobValid(blob *horizon.Blob) error {
	// blob should exists
	if blob == nil {
		return errBlobNotFound
	}

	// service expects specific blob type
	if blob.Type != kyc.KYCFormBlobType {
		return errInvalidBlobType
	}

	return nil
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
	return request.Details.KYC.AccountTypeToSet.Int == int(xdr.AccountTypeVerified)
}
