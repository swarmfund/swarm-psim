package marketmaker

import (
	"reflect"
	"time"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/regources"
)

type Config struct {
	CheckPeriod time.Duration `fig:"check_period"` // optional

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required"`

	AssetPairs []AssetPairConfig `fig:"asset_pairs,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"check_period": c.CheckPeriod.String(),
		"asset_pairs":  c.AssetPairs,
	}
}

type AssetPairConfig struct {
	BaseAsset        string           `fig:"base_asset,required"`
	QuoteAsset       string           `fig:"quote_asset,required"`
	BaseAssetVolume  regources.Amount `fig:"base_asset_volume"`
	QuoteAssetVolume regources.Amount `fig:"quote_asset_volume"`
	PriceMargin      float64          `fig:"price_margin,required"`
}

func (c AssetPairConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"base_asset":         c.BaseAsset,
		"quote_asset":        c.QuoteAsset,
		"base_asset_volume":  c.BaseAssetVolume,
		"quote_asset_volume": c.QuoteAssetVolume,
		"price_margin":       c.PriceMargin,
	}
}

var hooks = figure.Hooks{
	"[]marketmaker.AssetPairConfig": func(raw interface{}) (reflect.Value, error) {
		rawSlice, err := cast.ToSliceE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to cast to slice")
		}

		var configs []AssetPairConfig
		for _, rawElem := range rawSlice {
			rawConfig, err := cast.ToStringMapE(rawElem)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to cast slice element to map[string]interface{}", logan.F{
					"raw_element": rawElem,
				})
			}

			var config AssetPairConfig
			err = figure.
				Out(&config).
				From(rawConfig).
				With(figure.BaseHooks, utils.CommonHooks).
				Please()
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "Failed to figure out AssetPairConfig")
			}

			configs = append(configs, config)
		}

		return reflect.ValueOf(configs), nil
	},
}
