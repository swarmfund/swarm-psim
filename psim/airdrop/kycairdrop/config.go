package kycairdrop

import (
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	airdrop.IssuanceConfig `fig:"issuance,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`

	airdrop.EmailsConfig `fig:"emails,required"`

	BlackList []string `fig:"black_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"issuance":          c.IssuanceConfig,
		"emails":            c.EmailsConfig,
		"black_list_length": len(c.BlackList),
	}
}
