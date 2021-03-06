// Sequence diagram of this service can be found here:
// https://drive.google.com/file/d/12YisR3Pdf6jg4jTXKRDicZwUugO6B9DH/view?usp=sharing
package marketmaker

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/go/amount"
)

func init() {
	app.RegisterService(conf.ServiceMarketMaker, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config = Config{
		CheckPeriod: 30 * time.Second,
	}
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceMarketMaker)).
		With(figure.BaseHooks, utils.CommonHooks, hooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceMarketMaker,
		})
	}

	// TODO VALIDATE method
	for _, assetPair := range config.AssetPairs {
		if assetPair.PriceMargin <= 0 || assetPair.PriceMargin >= 100 * amount.One {
			return nil, errors.From(errors.New("PriceMargin must be bigger than zero and smaller than 100."), logan.F{
				"asset_pair": assetPair,
			})
		}

		if assetPair.BaseAssetVolume == 0 && assetPair.QuoteAssetVolume == 0 {
			return nil, errors.From(errors.New("BaseAssetVolume and QuoteAssetVolume cannot both be zero."), logan.F{
				"asset_pair": assetPair,
			})
		}
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	builder, err := horizonConnector.TXBuilder()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create TXBuilder")
	}

	return NewService(
		log,
		config,
		horizonConnector.Assets(),
		horizonConnector.Accounts(),
		horizonConnector.Submitter(),
		builder,
	), nil
}
