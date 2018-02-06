package resources

type Account struct {
	// TODO address
	AccountID              string                  `json:"account_id"`
	ExternalSystemAccounts []ExternalSystemAccount `json:"external_system_accounts"`
}

type ExternalSystemAccount struct {
	Type struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"type"`

	Address string `json:"data"`
}
