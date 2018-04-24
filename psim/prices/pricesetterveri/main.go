package pricesetterveri

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/prices/finder"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServicePriceSetterVerify, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServicePriceSetterVerify)).
		With(figure.BaseHooks, utils.ETHHooks, providers.FigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServicePriceSetterVerify,
		})
	}

	pFinder, err := newPriceFinder(ctx, log, config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create PriceFinder")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	listener, err := ape.Listener(config.Host, config.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init listener")
	}

	return New(
		conf.ServicePriceSetterVerify,
		log,
		config,
		pFinder,
		config.Signer,
		xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod),
		listener,
		globalConfig.Discovery(),
	), nil
}

func newPriceFinder(ctx context.Context, log *logan.Entry, config Config) (priceFinder, error) {
	// Set of PriceProviders
	usedProviders := map[string]struct{}{}
	var priceProviders []finder.PriceProvider

	quoteAsset := config.QuoteAsset
	if quoteAsset == "SUN" {
		quoteAsset = "USD"
	}

	for _, providerData := range config.Providers {
		if _, contains := usedProviders[providerData.Name]; contains {
			return nil, errors.From(errors.New("Duplication of PriceProviders not allowed."), logan.F{
				"provider_name": providerData.Name,
			})
		}

		usedProviders[providerData.Name] = struct{}{}
		specificProvider, err := providers.StartSpecificProvider(ctx, log, config.BaseAsset, quoteAsset,
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
