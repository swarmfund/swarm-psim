package resources

import "gitlab.com/tokend/horizon-connector/types"

type Balance struct {
	Asset     string       `json:"asset"`
	BalanceID string       `json:"balance_id"`
	AccountID string       `json:"account_id"`
	Balance   types.Amount `json:"balance"`
	Locked    types.Amount `json:"locked"`
}
