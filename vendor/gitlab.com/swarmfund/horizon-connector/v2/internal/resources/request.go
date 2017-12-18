package resources

import "gitlab.com/swarmfund/horizon-connector/v2/types"

type Request struct {
	ID          uint64         `json:"id,string"`
	PagingToken string         `json:"paging_token"`
	Hash        string         `json:"hash"`
	State       int32          `json:"request_state_i"`
	Details     RequestDetails `json:"details"`
}

type RequestDetails struct {
	RequestType int32                   `json:"request_type_i"`
	Withdraw    *RequestWithdrawDetails `json:"withdraw"`
}

type RequestWithdrawDetails struct {
	Amount           types.Amount `json:"amount"`
	BalanceID        string       `json:"balance_id"`
	DestinationAsset string       `json:"dest_asset_code"`
	ExternalDetails  string       `json:"external_details"`
}
