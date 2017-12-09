package horizon

type BalanceID struct {
	ID           string `json:"id,omitempty"`
	BalanceID    string `json:"balance_id"`
	AccountID    string `json:"account_id"`
	ExchangeID   string `json:"exchange_id"`
	ExchangeName string `json:"exchange_name"`
	Asset        string `json:"asset"`
}
