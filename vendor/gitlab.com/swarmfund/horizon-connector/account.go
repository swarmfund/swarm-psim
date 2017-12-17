package horizon

type Account struct {
	AccountID    string           `json:"account_id"`
	AccountType  int32            `json:"account_type_i"`
	Sequence     string           `json:"sequence"`
	BlockReasons int32  	      `json:"block_reasons"`
	Balances     []Balance        `json:"balances"`
	Signers      []Signer         `json:"signers"`
	Policies     AccountPolicies  `json:"policies"`
}

func (account *Account) BalanceByAsset(asset string) *Balance {
	for _, balance := range account.Balances {
		if balance.Asset == asset {
			return &balance
		}
	}
	return nil
}

type Balance struct {
	BalanceID      string `json:"balance_id"`
	AccountID      string `json:"account_id"`
	Balance        string `json:"balance"`
	ExchangeID     string `json:"exchange_id"`
	ExchangeName   string `json:"exchange_name"`
	Asset          string `json:"asset"`
	Locked         string `json:"locked"`
	RequireReview  bool   `json:"require_review"`
	StorageFee     string `json:"storage_fee"`
	StorageFeeTime string `json:"storage_fee_time"`
}

type AccountPolicies struct {
	Type int32 `json:"account_policies_type_i"`
}
