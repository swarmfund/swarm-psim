package ethwithdraw

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

const (
	// Here is the full list of RejectReasons, which Service can set into `reject_reason` of Request in case of validation error(s).
	RejectReasonMissingAddress    = "missing_address"
	RejectReasonAddressNotAString = "address_not_a_string"
	RejectReasonInvalidAddress    = "invalid_address"
	RejectReasonTooLittleAmount   = "too_little_amount"

	RequestStatePending  int32 = 1
	RequestStateApproved int32 = 3

	VersionPreConfirmDetailsKey = "version"
	TX1PreConfirmDetailsKey     = "raw_eth_tx_1"
	TX1HashPreConfirmDetailsKey = "eth_tx_1_hash"

	WithdrawAddressExtDetailsKey = "address"
)

// GetTX1 can return nil,nil if:
// - Request version is not 3, or
// - The Request is not of type TwoStepWithdraw
// in all other cases - nil error means non-nil Transaction and vice versa.
func getTX1(request horizon.Request) (*types.Transaction, error) {
	version, err := getPreConfirmationVersion(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get version of the Request")
	}
	// Change following if, if supported versions change.
	if version != 3 {
		// Not a v3 WithdrawRequests, this service is unable to process them.
		return nil, nil
	}

	if request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		return nil, nil
	}
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		return nil, errors.New("Unexpected RequestType, only TwoStepWithdraw(7) and Withdraw(4) are expected.")
	}

	preConfirmDetails := request.Details.Withdraw.PreConfirmationDetails

	rawTXHexI, ok := preConfirmDetails[TX1PreConfirmDetailsKey]
	if !ok {
		return nil, errors.New("Not found raw ETH TX_1 hex in the PreConfirmationDetails.")
	}
	rawTXHex, ok := rawTXHexI.(string)
	if !ok {
		return nil, errors.New("Raw ETH TX_1 in the PreConfirmationDetails is not of type string.")
	}
	fields := logan.F{
		"eth_tx_hex": rawTXHex,
	}

	tx, err := eth.Unmarshal(rawTXHex)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal ETH TX hex from PreConfirmationDetails", fields)
	}

	return tx, nil
}

func getPreConfirmationVersion(request horizon.Request) (float64, error) {
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		// Not a WithdrawRequest - either still TSWRequest or not a WithdrawalRequest at all.
		return 0, nil
	}

	versionI, ok := request.Details.Withdraw.PreConfirmationDetails[VersionPreConfirmDetailsKey]
	if !ok {
		// Not a v3 - old-style WithdrawRequests, this service is unable to process them.
		return 0, nil
	}
	version, ok := versionI.(float64)
	if !ok {
		return 0, errors.From(errors.New("Version in the ExternalDetails is not of type float64."), logan.F{
			"raw_version_value": versionI,
		})
	}
	return version, nil
}

func (s *Service) getTSWRejectReason(request horizon.Request, countedAssetAmount *big.Int) string {
	tswRequest := request.Details.TwoStepWithdraw

	addrI, ok := tswRequest.ExternalDetails[WithdrawAddressExtDetailsKey]
	if !ok {
		return RejectReasonMissingAddress
	}
	addr, ok := addrI.(string)
	if !ok {
		return RejectReasonAddressNotAString
	}
	if !common.IsHexAddress(addr) {
		return RejectReasonInvalidAddress
	}

	if countedAssetAmount.Cmp(s.config.MinWithdrawAmount) < -1 {
		return RejectReasonTooLittleAmount
	}

	return ""
}

func isProcessablePendingRequest(request horizon.Request) bool {
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		// Withdraw service only approves TwoStepWithdraw Requests, Withdraw Requests will be approved by Verify service.
		return false
	}
	if request.State != RequestStatePending {
		// We are only interested in pending TwoStepWithdrawRequests
		return false
	}

	return true
}
