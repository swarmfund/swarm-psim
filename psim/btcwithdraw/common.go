package btcwithdraw

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
)

// TODO Consider moving to some common package, as this logic is common for BTC and ETH.

const (
	// Here is the full list of RejectReasons, which Service can set into `reject_reason` of Request in case of validation error(s).
	RejectReasonInvalidAddress  RejectReason = "invalid_btc_address"
	RejectReasonTooLittleAmount RejectReason = "too_little_amount"
)

type RejectReason string

// GetWithdrawAddress obtains withdraw Address from the `address` field of the ExternalDetails
// of Withdraw in Request Details.
// Returns error if no `address` field in the ExternalDetails map or if the field is not a string.
func GetWithdrawAddress(request horizonV2.Request) (string, error) {
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

func GetWithdrawAmount(request horizonV2.Request) float64 {
	return float64(int64(request.Details.Withdraw.DestinationAmount)) / amount.One
}

// TODO Comment
func ValidateBTCAddress(addr string, defaultNet *chaincfg.Params) error {
	_, err := btcutil.DecodeAddress(addr, defaultNet)
	return err
}

// IsPendingBTCWithdraw returns true if the Request is of Withdraw type,
// is in pending state
// and its DestinationAsset is BTC.
func IsPendingBTCWithdraw(request horizonV2.Request) bool {
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		// not a withdraw request
		return false
	}

	if request.State != RequestStatePending {
		// State is not pending
		return false
	}

	if request.Details.Withdraw.DestinationAsset != BTCAsset {
		// Withdraw not to BTC.
		return false
	}

	return true
}

func GetRequestLoganFields(key string, request horizonV2.Request) logan.F {
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
