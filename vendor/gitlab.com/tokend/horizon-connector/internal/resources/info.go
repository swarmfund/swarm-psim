package resources

type Info struct {
	Passphrase         string `json:"network_passphrase"`
	MasterAccountID    string `json:"master_account_id"`
	TXExpirationPeriod int64  `json:"tx_expiration_period"`
}
