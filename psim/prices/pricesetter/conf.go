package pricesetter

import (
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/keypair"
	"reflect"
	"time"
)

type ProviderConfig struct {
	Name   string        `mapstructure:"name"`
	Period time.Duration `mapstructure:"period"`
}

type Config struct {
	BaseAsset            string           `mapstructure:"base_asset"`
	QuoteAsset           string           `mapstructure:"quote_asset"`
	Providers            []ProviderConfig `mapstructure:"providers"`
	ProvidersToAgree     int              `mapstructure:"providers_to_agree"`
	MaxPriceDeltaPercent string           `mapstructure:"max_price_delta_percent"`

	Source   keypair.Address `fig:"source"`
	SignerKP keypair.Full    `fig:"signer" mapstructure:"signer"`
}

var priceSetterFigureHooks = figure.Hooks{
	"[]pricesetter.ProviderConfig": func(raw interface{}) (reflect.Value, error) {
		providers, ok := raw.([]interface{})
		if !ok {
			return reflect.Value{}, errors.New("Unexpected type for providers")
		}

		result := make([]ProviderConfig, len(providers))
		for i := range providers {
			rawProvider, err := cast.ToStringMapE(providers[i])
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to cast provider to map[string]interface{}")
			}

			var provider ProviderConfig
			err = figure.
				Out(&provider).
				From(rawProvider).
				With(figure.BaseHooks).
				Please()
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to figure out provider")
			}
			result[i] = provider
		}

		return reflect.ValueOf(result), nil
	},
}
