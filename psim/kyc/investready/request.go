package investready

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

// ProveInterestingRequest returns non-nil error if the provided Request
// doesn't need to be considered by this Service.
func proveInterestingRequest(request horizon.Request) error {
	if request.State != kyc.RequestStatePending {
		// State is not pending
		return errors.Errorf("Invalid Request State (%d) expected Pending(%d).", request.State, kyc.RequestStatePending)
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
	if pendingTasks&(kyc.TaskSuperAdmin|kyc.TaskFaceValidation|kyc.TaskDocsExpirationDate|kyc.TaskSubmitIDMind|kyc.TaskCheckIDMind) != 0 {
		// Some checks hasn't been completed yet - too early to process this request.
		return errors.New("Some previous Tasks hasn't been approved yet - too early to process this Request.")
	}

	if pendingTasks&kyc.TaskCheckInvestReady == 0 {
		return errors.New("CheckInvestReady task is not set in pending tasks - ignoring this Request.")
	}

	return nil
}

func getInvestReadyUserHash(kycReq horizon.KYCRequest) string {
	var userHash string
	for _, extDetails := range kycReq.ExternalDetails {
		value, ok := extDetails[UserHashExtDetailsKey]
		if !ok {
			// No 'tx_id' key in these externalDetails.
			continue
		}

		userHash, ok = value.(string)
		if !ok {
			// UserHash field in ExternalDetails is not a string.

			// Must never happen, but just in case.
			// Maybe we need to log this shit here, if it happens..
			continue
		}
	}

	return userHash
}
