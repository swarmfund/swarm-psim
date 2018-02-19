package ratesync

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/ratesync/finder"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider/bitfinex"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider/bitstamp"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider/coinmarketcap"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider/gdax"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceRateSync, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceRateSync)).
		With(figure.BaseHooks, utils.ETHHooks, rateSyncFigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceRateSync,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.SignerKP)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	priceFinder, err := newPriceFinder(ctx, log, config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init PriceFinder")
	}

	return &service{
		baseAsset:  config.BaseAsset,
		quoteAsset: config.QuoteAsset,

		log:         log.WithField("runner", "price_setter"),
		source:      config.Source,
		signer:      config.SignerKP,
		connector:   horizonConnector.Submitter(),
		txBuilder:   xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod),
		priceFinder: priceFinder,
	}, nil
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
		specificProvider, err := startSpecificProvider(ctx, log, config, providerData)
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

func startSpecificProvider(ctx context.Context, log *logan.Entry, config Config, providerData Provider) (finder.PriceProvider, error) {
	switch providerData.Name {
	case bitfinex.Name:
		return priceProviderFromConnector(ctx, log, config, bitfinex.New(), providerData.Period), nil
	case bitstamp.Name:
		return priceProviderFromConnector(ctx, log, config, bitstamp.New(), providerData.Period), nil
	case coinmarketcap.Name:
		return priceProviderFromConnector(ctx, log, config, coinmarketcap.New(), providerData.Period), nil
	case gdax.Name:
		// Gdax exchange provides Prices over socket, so Connector which does htt.Get over ticker is unsuitable here.
		pointsStream, err := gdax.StartNewPriceStreamer(ctx, log, config.BaseAsset, config.QuoteAsset)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create and start Gdax Streamer")
		}

		return provider.StartNewProvider(ctx, providerData.Name, pointsStream, log), nil
	default:
		return nil, errors.From(errors.New("Unexpected PriceProvider name"), logan.F{
			"provider_name": providerData.Name,
		})
	}
}

func priceProviderFromConnector(
	ctx context.Context,
	log *logan.Entry,
	config Config,
	connector provider.Connector,
	period time.Duration) finder.PriceProvider {

	pointsStream := provider.StartNewPriceStreamer(ctx, log, config.BaseAsset, config.QuoteAsset, connector, period)
	return provider.StartNewProvider(ctx, connector.GetName(), pointsStream, log)
}
