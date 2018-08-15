package withdraw

import (
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

const (
	RequestStatePending int32 = 1
)

// TODO Functions in this file should probably become private.

var (
	ErrMissingWithdraw         = errors.New("Missing WithdrawRequest in the RequestDetails.")
	ErrMissingTwoStepWithdraw  = errors.New("Missing TwoStepWithdrawRequest in the RequestDetails.")
	ErrMissingRequestInDetails = errors.New("Missing both TwoStepWithdrawRequest and WithdrawRequest in the Request.")

	ErrMissingAddress    = errors.New("Missing address field in the ExternalDetails json.")
	ErrAddressNotAString = errors.New("Address field in ExternalDetails is not a string.")

	ErrMissingTXHex    = errors.New("Missing Offchain TX (tx_hex field) in the PreConfirmationDetails json of WithdrawRequest.")
	ErrTXHexNotAString = errors.New("Offchain TX (tx_hex field) in ExternalDetails of WithdrawRequest is not a string.")
)

// GetWithdrawalAddress obtains withdrawal Address from the `address` field of the ExternalDetails
// of Withdraw in Request Details.
//
// GetWithdrawalAddress does work well with both Withdraw and TwoStepWithdraw Requests.
//
// Returns error if no `address` field in the ExternalDetails map or if the field is not a string.
// Only returns errors with causes:
// - ErrMissingRequestInDetails
// - ErrMissingAddress
// - ErrAddressNotAString.
func GetWithdrawalAddress(request horizon.Request) (string, error) {
	if request.Details.TwoStepWithdraw != nil {
		return getWithdrawAddress(request.Details.TwoStepWithdraw.ExternalDetails)
	}

	if request.Details.Withdraw != nil {
		return getWithdrawAddress(request.Details.Withdraw.ExternalDetails)
	}

	return "", ErrMissingRequestInDetails
}

func getWithdrawAddress(externalDetails map[string]interface{}) (string, error) {
	addrValue, ok := externalDetails["address"]
	if !ok {
		return "", ErrMissingAddress
	}

	addr, ok := addrValue.(string)
	if !ok {
		return "", errors.From(ErrAddressNotAString, logan.F{"raw_address_value": addrValue})
	}

	return addr, nil
}

// TODO Comment
func GetWithdrawAmount(request regources.ReviewableRequest) (int64, error) {
	if request.Details.TwoStepWithdraw != nil {
		return int64(request.Details.TwoStepWithdraw.DestAssetAmount), nil
	}
	if request.Details.Withdraw != nil {
		return int64(request.Details.Withdraw.DestAssetAmount), nil
	}

	return 0, ErrMissingRequestInDetails
}

// GetTXHex obtains Withdraw TX hex from the `tx_hex` field of the ExternalDetails
// of Withdraw in Request Details.
//
// Returns error if Withdraw in Details is nil, or if no `tx_hex` field in the ExternalDetails map, or if the field is not a string.
// Only returns errors with causes equal to:
// - ErrMissingWithdraw
// - ErrMissingTXHex
// - ErrTXHexNotAString.
func GetTXHex(request regources.ReviewableRequest) (string, error) {
	if request.Details.Withdraw == nil {
		return "", ErrMissingWithdraw
	}

	txHexValue, ok := request.Details.Withdraw.PreConfirmationDetails["tx_hex"]
	if !ok {
		return "", ErrMissingTXHex
	}

	txHex, ok := txHexValue.(string)
	if !ok {
		return "", errors.From(ErrTXHexNotAString, logan.F{"raw_tx_hex_value": txHexValue})
	}

	return txHex, nil
}

// ProvePendingRequest returns empty string if the Request is:
// - in pending state;
// - type equals `neededRequestType`;
// - its DestinationAsset equals `asset`.
//
// Otherwise returns string describing the validation error.
func ProvePendingRequest(request regources.ReviewableRequest, asset string, neededRequestTypes ...int32) string {
	if request.State != RequestStatePending {
		// State is not pending
		return fmt.Sprintf("Invalid Request State (%d) expected Pending(%d).", request.State, RequestStatePending)
	}

	var isTypeAppropriate bool
	for _, neededRequestType := range neededRequestTypes {
		if request.Details.RequestType == neededRequestType {
			isTypeAppropriate = true
		}
	}
	if !isTypeAppropriate {
		return fmt.Sprintf("Invalid RequestType (%d) expected (%v).", request.Details.RequestType, neededRequestTypes)
	}

	var destAsset string
	if request.Details.TwoStepWithdraw != nil {
		destAsset = request.Details.TwoStepWithdraw.DestAssetCode
	}
	if request.Details.Withdraw != nil {
		destAsset = request.Details.Withdraw.DestAssetCode
	}
	// TODO If not Withdraw and not TSW - consider returning specific error (switch request.Details.RequestType)

	if destAsset != asset {
		// Withdraw not to BTC.
		return fmt.Sprintf("Wrong DestintationAsset (%s) expected BTC(%s).", destAsset, asset)
	}

	return ""
}
