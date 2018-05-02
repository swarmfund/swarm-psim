package providers

import (
	"reflect"
	"time"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ProviderConfig struct {
	Name   string        `mapstructure:"name,required"`
	Period time.Duration `mapstructure:"period,required"`
}

func (c ProviderConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"name":   c.Name,
		"period": c.Period.String(),
	}
}

var FigureHooks = figure.Hooks{
	"[]providers.ProviderConfig": func(raw interface{}) (reflect.Value, error) {
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
