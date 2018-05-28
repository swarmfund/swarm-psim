package resources

type Account struct {
	AccountID              string                  `json:"account_id"`
	IsBlocked              bool                    `json:"is_blocked"`
	AccountTypeI           int32                   `json:"account_type_i"`
	AccountType            string                  `json:"account_type"`
	ExternalSystemAccounts []ExternalSystemAccount `json:"external_system_accounts"`
	KYC                    AccountKYC              `json:"account_kyc"`
	Referrer               string                  `json:"referrer"`
}

type ExternalSystemAccount struct {
	Type struct {
		// Name human readable asset name
		Name string `json:"name"`
		// Value external system type
		Value int `json:"value"`
	} `json:"type"`
	// AssetCode TokenD asset code
	AssetCode string `json:"asset_code"`
	Address   string `json:"data"`
}

type AccountKYC struct {
	Data *KYCData `json:"KYCData"`
}

type KYCData struct {
	BlobID string `json:"blob_id"`
}
