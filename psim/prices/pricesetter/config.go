package pricesetter

import (
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	BaseAsset            string                     `mapstructure:"base_asset"`
	QuoteAsset           string                     `mapstructure:"quote_asset"`

	Providers            []providers.ProviderConfig `mapstructure:"providers"`

	ProvidersToAgree     int                        `mapstructure:"providers_to_agree"`
	MaxPriceDeltaPercent string                     `mapstructure:"max_price_delta_percent"`
	VerifierServiceName  string                     `fig:"verifier_service_name"`

	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`
}
