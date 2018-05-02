package pricesetter

import (
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	BaseAsset  string `mapstructure:"base_asset"`
	QuoteAsset string `mapstructure:"quote_asset"`

	Providers []providers.ProviderConfig `mapstructure:"providers"`

	ProvidersToAgree     int    `mapstructure:"providers_to_agree"`
	MaxPriceDeltaPercent string `mapstructure:"max_price_delta_percent"`
	VerifierServiceName  string `fig:"verifier_service_name"`

	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"base_asset":              c.BaseAsset,
		"quote_asset":             c.QuoteAsset,
		"providers":               c.Providers,
		"providers_to_agree":      c.ProvidersToAgree,
		"max_price_delta_percent": c.MaxPriceDeltaPercent,
		"verifier_service_name":   c.VerifierServiceName,
	}
}
