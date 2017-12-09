package create_account_streamer

// CreateAccountOp is used to store CreateAccount operation obtained from Horizon
type CreateAccountOp struct {
	Account string `json:"account"`
	// TODO Add BTC and ETH addresses here
}

type OperationsResponse struct {
	Embedded struct {
		Records []CreateAccountOp
	} `json:"_embedded"`

	Links struct {
		Next struct {
			HREF string `json:"href"`
		} `json:"next"`
	} `json:"_links"`
}
