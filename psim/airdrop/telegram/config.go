package telegram

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/listener"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Listener  listener.Config        `fig:"listener,required"`
	Issuance  airdrop.IssuanceConfig `fig:"issuance,required"`
	BlackList []string               `fig:"black_list"`

	TelegramSecretKey string          `json:"telegram_secret_key,required"`
	Source            keypair.Address `fig:"source,required"`
	Signer            keypair.Full    `fig:"signer,required" mapstructure:"signer"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"listener":          c.Listener,
		"issuance":          c.Issuance,
		"black_list_length": len(c.BlackList),
	}
}

func (c Config) Validate() (validationErr error) {
	// Wrap of nil error always returns nil.
	return errors.Wrap(c.Issuance.Validate(), "Issuance config is invalid")
}
