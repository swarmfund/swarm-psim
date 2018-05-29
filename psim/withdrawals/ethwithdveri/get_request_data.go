package ethwithdveri

import (
	"math/big"

	"fmt"

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

	TX1HashPreConfirmDetailsKey = "eth_tx_1_hash"
	TX2ReviewerDetailsKey       = "raw_eth_tx_2"
	TX2HashReviewerDetailsKey   = "eth_tx_2_hash"

	WithdrawAddressExtDetailsKey = "address"
	VersionExtDetailsKey         = "version"
)

// GetTX2 can return nil,nil if:
// - WithdrawRequest version is not 2, or
// - The Request is not of type Withdraw
// - WithdrawRequest is not approved
// in all other cases - nil error means non-nil Transaction and vice versa.
func getTX2(request horizon.Request) (string, *types.Transaction, error) {
	version, err := getVersion(request)
	if err != nil {
		return "", nil, errors.Wrap(err, "Failed to get version of the Request")
	}
	// Change following if, if supported versions change.
	if version != 2 {
		// Not a v2 WithdrawRequests, this service is unable to process them.
		return "", nil, nil
	}

	if request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		return "", nil, nil
	}
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		return "", nil, errors.New("Unexpected RequestType, only TwoStepWithdraw(7) and Withdraw(4) are expected.")
	}

	if request.State != RequestStateApproved {
		// We are looking for TX2 only in approved WithdrawRequests
		return "", nil, nil
	}

	reviewerDetails := request.Details.Withdraw.ReviewerDetails

	rawTXHexI, ok := reviewerDetails[TX2ReviewerDetailsKey]
	if !ok {
		return "", nil, errors.New("Not found raw ETH TX_2 hex in the ReviewerDetails.")
	}
	rawTXHex, ok := rawTXHexI.(string)
	if !ok {
		return "", nil, errors.New("Raw ETH TX_2 in the ReviewerDetails is not of type string.")
	}
	fields := logan.F{
		"eth_tx_hex": rawTXHex,
	}

	tx, err := eth.Unmarshal(rawTXHex)
	if err != nil {
		return "", nil, errors.Wrap(err, "Failed to unmarshal ETH TX hex from ReviewerDetails", fields)
	}

	return rawTXHex, tx, nil
}

// TODO Avoid duplication with ethwithdraw service.
func getVersion(request horizon.Request) (float64, error) {
	var extDetails map[string]interface{}

	switch request.Details.RequestType {
	case int32(xdr.ReviewableRequestTypeWithdraw):
		extDetails = request.Details.Withdraw.ExternalDetails
	case int32(xdr.ReviewableRequestTypeTwoStepWithdrawal):
		extDetails = request.Details.TwoStepWithdraw.ExternalDetails
	default:
		return 0, errors.New("Unexpected RequestType, only TwoStepWithdraw(7) and Withdraw(4) are expected.")
	}

	versionI, ok := extDetails[VersionExtDetailsKey]
	if !ok {
		// Not a v2 - old-style WithdrawRequests, this service is unable to process them.
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

func (s *Service) getWithdrawRejectReason(request horizon.Request, countedAssetAmount *big.Int) string {
	withdrawRequest := request.Details.Withdraw

	addrI, ok := withdrawRequest.ExternalDetails[WithdrawAddressExtDetailsKey]
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

func getRequestNotProcessableReason(request horizon.Request) string {
	// Just in case, should never happen as filters only request Withdrawal requests
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) && request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		return fmt.Sprintf("Invalid RequestType (%d).", request.Details.RequestType)
	}

	if request.State != RequestStatePending {
		// We are only interested in pending WithdrawRequests
		return fmt.Sprintf("Invalid Request State (%d).", request.State)

	}

	return ""
}
