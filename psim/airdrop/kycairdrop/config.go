package kycairdrop

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	IssuanceConfig   airdrop.IssuanceConfig `fig:"issuance,required"`
	USACheckDisabled bool                   `fig:"usa_check_disabled"`
	Source           keypair.Address        `fig:"source,required"`
	Signer           keypair.Full           `fig:"signer,required" mapstructure:"signer"`
	EmailsConfig     airdrop.EmailsConfig   `fig:"emails,required"`
	BlackList        []string               `fig:"black_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"issuance":          c.IssuanceConfig,
		"emails":            c.EmailsConfig,
		"black_list_length": len(c.BlackList),
	}
}

func (c Config) Validate() error {
	return errors.Wrap(c.EmailsConfig.Validate(), "config is invalid")
}
