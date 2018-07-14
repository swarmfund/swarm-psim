package pricesetterveri

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Host string `fig:"host"`
	Port int    `fig:"port"`

	BaseAsset            string                     `fig:"base_asset,required"`
	QuoteAsset           string                     `fig:"quote_asset,required"`
	Providers            []providers.ProviderConfig `fig:"providers,required"`
	ProvidersToAgree     int                        `fig:"providers_to_agree,required"`
	MaxPriceDeltaPercent string                     `fig:"max_price_delta_percent,required"`

	Signer              keypair.Full `fig:"signer,required"`
	VerifierServiceName string       `fig:"verifier_service_name,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"host": c.Host,
		"port": c.Port,

		"base_asset": c.BaseAsset,
		"quote_asset": c.QuoteAsset,
		"providers": c.Providers,
		"providers_to_agree": c.ProvidersToAgree,
		"max_price_delta_percent": c.MaxPriceDeltaPercent,

		"verifier_service_name": c.VerifierServiceName,
	}
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{}
	err := figure.
		Out(config).
		From(configData).
		With(figure.BaseHooks, utils.ETHHooks, providers.FigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")
	}

	return config, nil
}
