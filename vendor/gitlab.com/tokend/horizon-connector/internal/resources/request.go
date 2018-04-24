package resources

import (
	"gitlab.com/tokend/horizon-connector/types"
	"time"
)

type Request struct {
	ID           uint64         `json:"id,string"`
	PagingToken  string         `json:"paging_token"`
	Hash         string         `json:"hash"`
	RejectReason string         `json:"reject_reason"`
	State        int32          `json:"request_state_i"`
	Details      RequestDetails `json:"details"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (r Request) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"id":              r.ID,
		"paging_token":    r.PagingToken,
		"hash":            r.Hash,
		"request_state_i": r.State,
		"details":         r.Details,
		"created_at":      r.CreatedAt,
		"updated_at":      r.UpdatedAt,
		"reject_reason":   r.RejectReason,
	}
}

type RequestDetails struct {
	RequestType     int32                   `json:"request_type_i"`
	Withdraw        *RequestWithdrawDetails `json:"withdraw"`
	TwoStepWithdraw *RequestWithdrawDetails `json:"two_step_withdrawal"`
	KYC             *RequestKYCDetails      `json:"update_kyc"`
}

func (d RequestDetails) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"request_type_i":      d.RequestType,
		"withdraw":            d.Withdraw,
		"two_step_withdrawal": d.TwoStepWithdraw,
		"kyc":                 d.KYC,
	}
}

type RequestWithdrawDetails struct {
	Amount            types.Amount           `json:"amount"`
	DestinationAmount types.Amount           `json:"dest_asset_amount"`
	DestinationAsset  string                 `json:"dest_asset_code"`
	BalanceID         string                 `json:"balance_id"`
	ExternalDetails   map[string]interface{} `json:"external_details"`

	FixedFee               types.Amount           `json:"fixed_fee"`
	PercentFee             types.Amount           `json:"percent_fee"`
	PreConfirmationDetails map[string]interface{} `json:"pre_confirmation_details"`
	ReviewerDetails        map[string]interface{} `json:"reviewer_details"`
}

func (d RequestWithdrawDetails) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"amount":             d.Amount,
		"destination_amount": d.DestinationAmount,
		"destination_asset":  d.DestinationAsset,
		"balance_id":         d.BalanceID,
		"external_details":   d.ExternalDetails,

		"fixed_fee":                d.FixedFee,
		"percent_fee":              d.PercentFee,
		"pre_confirmation_details": d.PreConfirmationDetails,
		"reviewer_details":         d.ReviewerDetails,
	}
}

type RequestKYCDetails struct {
	AccountToUpdateKYC string                   `json:"account_to_update_kyc"`
	AccountTypeToSet   AccountTypeToSet         `json:"account_type_to_set"`
	KYCLevel           uint32                   `json:"kyc_level"`
	KYCData            map[string]interface{}   `json:"kyc_data"`
	AllTasks           uint32                   `json:"all_tasks"`
	PendingTasks       uint32                   `json:"pending_tasks"`
	SequenceNumber     uint32                   `json:"sequence_number"`
	ExternalDetails    []map[string]interface{} `json:"external_details"`
}

func (d RequestKYCDetails) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_to_update_kyc": d.AccountToUpdateKYC,
		"account_type_to_set":   d.AccountTypeToSet.String,
		"kyc_level_i":           d.KYCLevel,
		"kyc_data_map":          d.KYCData,
		"all_tasks":             d.AllTasks,
		"pending_tasks":         d.PendingTasks,
		"sequence_number":       d.SequenceNumber,
		"external_details":      d.ExternalDetails,
	}
}

type AccountTypeToSet struct {
	Int    int    `json:"int"`
	String string `json:"string"`
}
