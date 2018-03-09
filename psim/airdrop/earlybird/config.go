package earlybird

import (
	"time"

	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Asset            string     `fig:"issuance_asset"`
	Amount           uint64     `fig:"issuance_amount"`
	RegisteredBefore *time.Time `fig:"registered_before"`

	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	airdrop.EmailsConfig `fig:"emails"`

	WhiteList []string `fig:"white_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"issuance_asset":    c.Asset,
		"issuance_amount":   c.Amount,
		"registered_before": c.RegisteredBefore.String(),
		"emails":            c.EmailsConfig,
		"white_list_len":    len(c.WhiteList),
	}
}
