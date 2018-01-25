package withdraw

import (
	"encoding/json"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

const (
	RequestStatePending int32 = 1
)

// TODO Functions in this file should probably become private.

var (
	ErrMissingWithdraw        = errors.New("Missing  in the ExternalDetails json of WithdrawRequest.")
	ErrMissingTwoStepWithdraw = errors.New("Missing  in the ExternalDetails json of TowStepWithdrawRequest.")

	ErrMissingAddress    = errors.New("Missing address field in the ExternalDetails json of TwoStepWithdrawRequest.")
	ErrAddressNotAString = errors.New("Address field in ExternalDetails of WithdrawalRequest is not a string.")

	ErrMissingTXHex    = errors.New("Missing Offchain TX (tx_hex field) in the PreConfirmationDetails json of WithdrawRequest.")
	ErrTXHexNotAString = errors.New("Offchain TX (tx_hex field) in ExternalDetails of WithdrawRequest is not a string.")
)

// GetWithdrawAddress obtains withdraw Address from the `address` field of the ExternalDetails
// of Withdraw in Request Details.
//
// Returns error if no `address` field in the ExternalDetails map or if the field is not a string.
// Only returns errors with causes:
// - ErrMissingTwoStepWithdraw
// - ErrMissingAddress
// - ErrAddressNotAString.
func GetWithdrawAddress(request horizon.Request) (string, error) {
	if request.Details.TwoStepWithdraw == nil {
		return "", ErrMissingTwoStepWithdraw
	}

	addrValue, ok := request.Details.TwoStepWithdraw.ExternalDetails["address"]
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
func GetWithdrawAmount(request horizon.Request) (int64, error) {
	if request.Details.TwoStepWithdraw != nil {
		return int64(request.Details.TwoStepWithdraw.DestinationAmount), nil
	}
	if request.Details.Withdraw != nil {
		return int64(request.Details.Withdraw.DestinationAmount), nil
	}

	return 0, ErrMissingTwoStepWithdraw
}

// GetTXHex obtains Withdraw TX hex from the `tx_hex` field of the ExternalDetails
// of Withdraw in Request Details.
//
// Returns error if no `tx_hex` field in the ExternalDetails map or if the field is not a string.
// Only returns errors with causes equal to either ErrMissingTXHex or ErrTXHexNotAString.
func GetTXHex(request horizon.Request) (string, error) {
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

// ObtainRequest gets Request by the requestID from Horizon, using provided horizonConnector.
// Error means that could not get response from Horizon
// or failed to unmarshal the response into horizon.Request.
func ObtainRequest(horizonClient *horizon.Client, requestID uint64) (*horizon.Request, error) {
	respBytes, err := horizonClient.Get(fmt.Sprintf("/requests/%d", requestID))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Request from Horizon")
	}

	var request horizon.Request
	err = json.Unmarshal(respBytes, &request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal Request from the Horizon response", logan.F{
			"horizon_response": string(respBytes),
		})
	}

	return &request, nil
}

// ProvePendingRequest returns empty string if the Request is:
// - in pending state;
// - type equals `neededRequestType`;
// - its DestinationAsset equals `asset`.
//
// Otherwise returns string describing the validation error.
func ProvePendingRequest(request horizon.Request, asset string, neededRequestTypes ...int32) string {
	if request.State != RequestStatePending {
		// State is not pending
		return fmt.Sprintf("Invalid Request State (%d) expected Pending(%d).", request.State, RequestStatePending)
	}

	var isTypeValid bool
	for _, neededRequestType := range neededRequestTypes {
		if request.Details.RequestType == neededRequestType {
			//return fmt.Sprintf("Invalid RequestType (%d) expected (%d).", request.Details.RequestType, neededRequestType)
			isTypeValid = true
		}
	}
	if !isTypeValid {
		return fmt.Sprintf("Invalid RequestType (%d) expected (%v).", request.Details.RequestType, neededRequestTypes)
	}

	var destAsset string
	if request.Details.TwoStepWithdraw != nil {
		destAsset = request.Details.TwoStepWithdraw.DestinationAsset
	}
	if request.Details.Withdraw != nil {
		destAsset = request.Details.Withdraw.DestinationAsset
	}

	if destAsset != asset {
		// Withdraw not to BTC.
		return fmt.Sprintf("Wrong DestintationAsset (%s) expected BTC(%s).", destAsset, asset)
	}

	return ""
}
