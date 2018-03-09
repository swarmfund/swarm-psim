package kycairdrop

import (
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Asset  string `fig:"issuance_asset"`
	Amount uint64 `fig:"issuance_amount"`

	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	airdrop.EmailsConfig `fig:"emails"`

	BlackList []string `fig:"black_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"asset":             c.Asset,
		"amount":            c.Amount,
		"emails":            c.EmailsConfig,
		"black_list_length": len(c.BlackList),
	}
}
