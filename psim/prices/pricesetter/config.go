package pricesetter

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	BaseAsset  string `fig:"base_asset,required" mapstructure:"base_asset,required"`
	QuoteAsset string `fig:"quote_asset,required" mapstructure:"quote_asset,required"`

	Providers []providers.ProviderConfig `fig:"providers,required" mapstructure:"providers,required"`

	SubmitPeriod         time.Duration `fig:"submit_period,required"`
	ProvidersToAgree     int           `fig:"providers_to_agree,required" mapstructure:"providers_to_agree,required"`
	MaxPriceDeltaPercent string        `fig:"max_price_delta_percent,required" mapstructure:"max_price_delta_percent,required"`
	// DisableVerify if true service will not seek verification and submit price update on it's own
	DisableVerify bool `fig:"disable_verify"`
	// VerifierServiceName discovery service name which service will seek verification from.
	// Required if DisableVerify is false
	VerifierServiceName  string        `fig:"verifier_service_name"`

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

	if !c.DisableVerify && c.VerifierServiceName == "" {
		return errors.New("verifier_service_name is required")
	}

	return nil
}
