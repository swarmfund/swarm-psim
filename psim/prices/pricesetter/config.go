package pricesetter

import (
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	BaseAsset  string `mapstructure:"base_asset,required"`
	QuoteAsset string `mapstructure:"quote_asset,required"`

	Providers []providers.ProviderConfig `mapstructure:"providers,required"`

	ProvidersToAgree     int    `mapstructure:"providers_to_agree,required"`
	MaxPriceDeltaPercent string `mapstructure:"max_price_delta_percent,required"`
	VerifierServiceName  string `fig:"verifier_service_name,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer,required"`
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
