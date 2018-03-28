package idmind

import "gitlab.com/tokend/keypair"

type Config struct {
	Connector ConnectorConfig `fig:"connector,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`

	WhiteList []string `fig:"white_list"`
	BlackList []string `fig:"black_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"connector": c.Connector,

		"white_list_len": len(c.WhiteList),
		"black_list_len": len(c.BlackList),
	}
}
