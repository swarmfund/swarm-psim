package ratesync

import (
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v2/errors"
	"gitlab.com/tokend/keypair"
	"reflect"
	"time"
)

type Provider struct {
	Name   string        `mapstructure:"name"`
	Period time.Duration `mapstructure:"period"`
}

type Config struct {
	BaseAsset            string     `mapstructure:"base_asset"`
	QuoteAsset           string     `mapstructure:"quote_asset"`
	Providers            []Provider `mapstructure:"providers"`
	ProvidersToAgree     int        `mapstructure:"providers_to_agree"`
	MaxPriceDeltaPercent string     `mapstructure:"max_price_delta_percent"`

	Source   keypair.Address `fig:"source"`
	SignerKP keypair.Full    `fig:"signer" mapstructure:"signer"`
}

var rateSyncFigureHooks = figure.Hooks{
	"[]ratesync.Provider": func(raw interface{}) (reflect.Value, error) {
		providers, ok := raw.([]interface{})
		if !ok {
			return reflect.Value{}, errors.New("Unexpected type for providers")
		}

		result := make([]Provider, len(providers))
		for i := range providers {
			rawProvider, err := cast.ToStringMapE(providers[i])
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to cast provider to map[string]interface{}")
			}

			var provider Provider
			err = figure.
				Out(&provider).
				From(rawProvider).
				With(figure.BaseHooks).
				Please()
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to figureout provider")
			}
			result[i] = provider
		}

		return reflect.ValueOf(result), nil
	},
}
