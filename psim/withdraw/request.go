package withdraw

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

const (
	RequestStatePending int32 = 1
)

var (
	ErrMissingAddress    = errors.New("Missing field in the ExternalDetails json of WithdrawalRequest.")
	ErrAddressNotAString = errors.New("Address field in ExternalDetails of WithdrawalRequest is not a string.")
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
// Returns error if no `address` field in the ExternalDetails map or if the field is not a string.
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
