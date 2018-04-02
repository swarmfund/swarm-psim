package resources

type Account struct {
	// TODO address
	AccountID              string                  `json:"account_id"`
	ExternalSystemAccounts []ExternalSystemAccount `json:"external_system_accounts"`
	KYC AccountKYC `json:"account_kyc"`
}

type ExternalSystemAccount struct {
	Type struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"type"`

	Address string `json:"data"`
}

type AccountKYC struct {
	Data *KYCData `json:"KYCData"`
}

type KYCData struct {
	BlobID string `json:"blob_id"`
}
