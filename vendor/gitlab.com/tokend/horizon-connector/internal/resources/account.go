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

func (a Account) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id":               a.AccountID,
		"is_blocked":               a.IsBlocked,
		"account_type_i":           a.AccountTypeI,
		"account_type":             a.AccountType,
		"external_system_accounts": a.ExternalSystemAccounts,
		"kyc":      a.KYC,
		"referrer": a.Referrer,
	}
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

func (k AccountKYC) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"data": k.Data,
	}
}

type KYCData struct {
	BlobID string `json:"blob_id"`
}

func (d KYCData) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"blob_id": d.BlobID,
	}
}
