package withdraw

import (
	"encoding/json"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

const (
	RequestStatePending int32 = 1
)

var (
	ErrMissingAddress    = errors.New("Missing address field in the ExternalDetails json of WithdrawalRequest.")
	ErrMissingTXHex      = errors.New("Missing Offchain TX (tx_hex field) in the ExternalDetails json of WithdrawalRequest.")
	ErrAddressNotAString = errors.New("Address field in ExternalDetails of WithdrawalRequest is not a string.")
	ErrTXHexNotAString   = errors.New("Offchain TX (tx_hex field) in ExternalDetails of WithdrawalRequest is not a string.")
)

// GetRequestLoganFields is a helper which builds map of logan.F for logging, so that not to do this
// each time horizon.Request needs to be logged.
//
// This method exists because of the lack of GetLoganFields() method on the horizon.Request type.
func GetRequestLoganFields(key string, request horizon.Request) logan.F {
	result := logan.F{
		key + "_id":    request.ID,
		key + "_state": request.State,
	}

	if request.Details.Withdraw != nil {
		detKey := key + "_withdraw_details"

		result[detKey+"_amount"] = request.Details.Withdraw.Amount
		result[detKey+"_destination_amount"] = request.Details.Withdraw.DestinationAmount
		result[detKey+"_balance_id"] = request.Details.Withdraw.BalanceID
		result[detKey+"_external_details"] = request.Details.Withdraw.ExternalDetails
	}

	return result
}

// GetWithdrawAddress obtains withdraw Address from the `address` field of the ExternalDetails
// of Withdraw in Request Details.
//
// Returns error if no `address` field in the ExternalDetails map or if the field is not a string.
// Only returns errors with causes either ErrMissingAddress or ErrAddressNotAString.
func GetWithdrawAddress(request horizon.Request) (string, error) {
	addrValue, ok := request.Details.Withdraw.ExternalDetails["address"]
	if !ok {
		return "", ErrMissingAddress
	}

	addr, ok := addrValue.(string)
	if !ok {
		return "", errors.From(ErrAddressNotAString, logan.F{"raw_address_value": addrValue})
	}

	return addr, nil
}

// GetWithdrawAmount retrieves DestinationAmount of the Withdraw from Details of the Request
// and divides this value by the amount.One (the value of one whole unit of currency).
func GetWithdrawAmount(request horizon.Request) float64 {
	return float64(int64(request.Details.Withdraw.DestinationAmount)) / amount.One
}

// GetTXHex obtains Withdraw TX hex from the `tx_hex` field of the ExternalDetails
// of Withdraw in Request Details.
//
// Returns error if no `tx_hex` field in the ExternalDetails map or if the field is not a string.
// Only returns errors with causes equal to either ErrMissingTXHex or ErrTXHexNotAString.
func GetTXHex(request horizon.Request) (string, error) {
	txHexValue, ok := request.Details.Withdraw.ExternalDetails["tx_hex"]
	if !ok {
		return "", ErrMissingTXHex
	}

	txHex, ok := txHexValue.(string)
	if !ok {
		return "", errors.From(ErrTXHexNotAString, logan.F{"raw_tx_hex_value": txHexValue})
	}

	return txHex, nil
}

// TODO Comment
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
func ProvePendingRequest(request horizon.Request, neededRequestType *int32, asset string) string {
	if request.State != RequestStatePending {
		// State is not pending
		return fmt.Sprintf("Invalid Request State (%d) expected Pending(%d).", request.State, RequestStatePending)
	}

	if neededRequestType != nil {
		if request.Details.RequestType != *neededRequestType {
			// not a withdraw request
			return fmt.Sprintf("Invalid RequestType (%d) expected (%d).", request.Details.RequestType, neededRequestType)
		}
	}

	if request.Details.Withdraw.DestinationAsset != asset {
		// Withdraw not to BTC.
		return fmt.Sprintf("Wrong DestintationAsset (%s) expected BTC(%s).", request.Details.Withdraw.DestinationAsset, asset)
	}

	return ""
}
