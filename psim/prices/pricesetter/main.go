package pricesetter

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/prices/finder"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServicePriceSetter, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServicePriceSetter)).
		With(figure.BaseHooks, utils.ETHHooks, providers.FigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServicePriceSetter,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	priceFinder, err := newPriceFinder(ctx, log, config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init PriceFinder")
	}

	txBuilder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return newService(config, log, horizonConnector.Submitter(), priceFinder, txBuilder), nil
}

func newPriceFinder(ctx context.Context, log *logan.Entry, config Config) (priceFinder, error) {
	// Set of PriceProviders
	usedProviders := map[string]struct{}{}
	var priceProviders []finder.PriceProvider

	for _, providerData := range config.Providers {
		if _, contains := usedProviders[providerData.Name]; contains {
			return nil, errors.From(errors.New("Duplication of PriceProviders not allowed."), logan.F{
				"provider_name": providerData.Name,
			})
		}

		usedProviders[providerData.Name] = struct{}{}
		specificProvider, err := providers.StartSpecificProvider(ctx, log, config.BaseAsset, config.QuoteAsset,
			providerData.Name, providerData.Period)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to init specific PriceProvider")
		}

		priceProviders = append(priceProviders, specificProvider)
	}

	maxPriceDelta, err := amount.Parse(config.MaxPriceDeltaPercent)
	if err != nil {
		return nil, errors.From(errors.New("Failed to parse maxPriceDelta."), logan.F{
			"raw_max_price_delta_percent_from_config": config.MaxPriceDeltaPercent,
		})
	}

	priceFinder, err := finder.NewPriceFinder(log, priceProviders, maxPriceDelta, config.ProvidersToAgree)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init price finder")
	}

	return priceFinder, nil
}
