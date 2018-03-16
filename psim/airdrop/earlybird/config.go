package earlybird

import (
	"time"

	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	airdrop.IssuanceConfig `fig:"issuance,required"`

	RegisteredBefore *time.Time `fig:"registered_before,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`

	airdrop.EmailsConfig `fig:"emails,required"`

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
