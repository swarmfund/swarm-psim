package pricesetterveri

import (
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Host string `fig:"host"`
	Port int    `fig:"port"`

	BaseAsset            string                     `mapstructure:"base_asset"`
	QuoteAsset           string                     `mapstructure:"quote_asset"`
	Providers            []providers.ProviderConfig `mapstructure:"providers"`
	ProvidersToAgree     int                        `mapstructure:"providers_to_agree"`
	MaxPriceDeltaPercent string                     `mapstructure:"max_price_delta_percent"`

	Signer keypair.Full `fig:"signer" mapstructure:"signer"`
}
