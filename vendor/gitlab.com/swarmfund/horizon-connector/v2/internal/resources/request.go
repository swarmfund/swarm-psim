package resources

import (
	"gitlab.com/swarmfund/horizon-connector/v2/types"
)

type Request struct {
	ID          uint64         `json:"id,string"`
	PagingToken string         `json:"paging_token"`
	Hash        string         `json:"hash"`
	State       int32          `json:"request_state_i"`
	Details     RequestDetails `json:"details"`
}

func (r Request) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"id":              r.ID,
		"paging_token":    r.PagingToken,
		"hash":            r.Hash,
		"request_state_i": r.State,
		"details":         r.Details,
	}
}

type RequestDetails struct {
	RequestType     int32                   `json:"request_type_i"`
	Withdraw        *RequestWithdrawDetails `json:"withdraw"`
	TwoStepWithdraw *RequestWithdrawDetails `json:"two_step_withdrawal"`
}

func (d RequestDetails) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"request_type_i":      d.RequestType,
		"withdraw":            d.Withdraw,
		"two_step_withdrawal": d.TwoStepWithdraw,
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
