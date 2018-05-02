package pricesetter

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	BaseAsset  string `mapstructure:"base_asset,required"`
	QuoteAsset string `mapstructure:"quote_asset,required"`

	Providers []providers.ProviderConfig `mapstructure:"providers,required"`

	SubmitPeriod         time.Duration `fig:"submit_period,required"`
	ProvidersToAgree     int           `mapstructure:"providers_to_agree,required"`
	MaxPriceDeltaPercent string        `mapstructure:"max_price_delta_percent,required"`
	VerifierServiceName  string        `fig:"verifier_service_name,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"base_asset":              c.BaseAsset,
		"quote_asset":             c.QuoteAsset,
		"providers":               c.Providers,
		"submit_period":           c.SubmitPeriod,
		"providers_to_agree":      c.ProvidersToAgree,
		"max_price_delta_percent": c.MaxPriceDeltaPercent,
		"verifier_service_name":   c.VerifierServiceName,
	}
}

func (c Config) Validate() (validationErr error) {
	var maxPeriodProvider providers.ProviderConfig

	for _, provider := range c.Providers {
		if provider.Period > maxPeriodProvider.Period {
			maxPeriodProvider = provider
		}
	}

	if c.SubmitPeriod <= maxPeriodProvider.Period {
		return errors.From(errors.New("SubmitPeriod must be grater than the max period among all Providers."), logan.F{
			"max_period_provider": maxPeriodProvider,
			"submit_period":       c.SubmitPeriod,
		})
	}

	return nil
}
