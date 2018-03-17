package ordernotifier

import (
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	EmailsConfig `fig:"emails"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"emails": c.EmailsConfig,
	}
}
