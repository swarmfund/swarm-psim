package horizon

import (
	"encoding/json"
	"net/http"
)

type Info struct {
	NetworkPassphrase    string `json:"network_passphrase"`
	CommissionAccountID  string `json:"commission_account_id"`
	MasterAccountID      string `json:"master_account_id"`
	OperationalAccountID string `json:"operational_account_id"`
	StorageFeeAccountID  string `json:"storage_fee_account_id"`
	TxExpirationPeriod   int64  `json:"tx_expiration_period"`
}

func NewInfo(horizon string) (*Info, error) {
	resp, err := http.Get(horizon)
	if err != nil {
		return nil, err
	}

	defer func() { err = resp.Body.Close() }()

	var info Info
	err = json.NewDecoder(resp.Body).Decode(&info)

	return &info, err
}
