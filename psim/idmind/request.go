package idmind

import (
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	RequestStatePending int32 = 1
)

// TODO
func proveInterestingRequest(request horizon.Request) error {
	if request.State != RequestStatePending {
		// State is not pending
		return errors.Errorf("Invalid Request State (%d) expected Pending(%d).", request.State, RequestStatePending)
	}

	//details := request.Details
	//
	//if details.RequestType != xdr.ReviewableRequestTypeKYC {
	//	return errors.Errorf("Invalid RequestType (%d) expected KYC(%d).", details.RequestType, xdr.ReviewableRequestTypeKYC)
	//}

	return nil
}
