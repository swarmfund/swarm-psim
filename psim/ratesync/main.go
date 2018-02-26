package ratesync

import (
	"context"

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
	"time"
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
		return nil, errors.Wrap(err, "failed to init price finder")
	}

	return &service{
		baseAsset:   config.BaseAsset,
		quoteAsset:  config.QuoteAsset,
		log:         log.WithField("runner", "ratesync"),
		source:      config.Source,
		signer:      config.SignerKP,
		connector:   horizonConnector.Submitter(),
		builder:     xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod),
		priceFinder: priceFinder,
	}, nil
}

func newPriceFinder(ctx context.Context, log *logan.Entry, config Config) (priceFinder, error) {
	usedProviders := map[string]bool{}
	var ratesProviders []finder.RatesProvider
	for _, providerData := range config.Providers {
		if _, contains := usedProviders[providerData.Name]; contains {
			return nil, errors.New("Duplication of providers not allowed: " + providerData.Name)
		}

		usedProviders[providerData.Name] = true
		specificProvider, err := getSpecificProvider(ctx, log, config, providerData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init specific price provider")
		}

		ratesProviders = append(ratesProviders, specificProvider)
	}

	maxPriceDelta, err := amount.Parse(config.MaxPriceDeltaPercent)
	if err != nil {
		return nil, errors.New("failed to parse max price delta " + config.MaxPriceDeltaPercent)
	}

	priceFinder, err := finder.NewPriceFinder(log, ratesProviders, maxPriceDelta, config.ProvidersToAgree)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init price finder")
	}

	return priceFinder, nil
}

func getSpecificProvider(ctx context.Context, log *logan.Entry, config Config, providerData Provider) (finder.RatesProvider, error) {
	switch providerData.Name {
	case "bitfinex":
		return providerFromConnector(ctx, log, config, bitfinex.New(), providerData.Period), nil
	case "bitstamp":
		return providerFromConnector(ctx, log, config, bitstamp.New(), providerData.Period), nil
	case "coinmarketcap":
		return providerFromConnector(ctx, log, config, coinmarketcap.New(), providerData.Period), nil
	case "gdax":
		streamer, err := gdax.StartNewPriceStreamer(ctx, log, config.BaseAsset, config.QuoteAsset)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create gdax steamer")
		}

		return provider.StartNewProvider(ctx, providerData.Name, streamer, log), nil
	default:
		return nil, errors.New("Unexpected provider: " + providerData.Name)
	}
}

func providerFromConnector(ctx context.Context, log *logan.Entry,
	config Config, connector provider.Connector, period time.Duration) finder.RatesProvider {

	streamer := provider.StartNewPriceStreamer(ctx, log, config.BaseAsset, config.QuoteAsset, connector, period)
	return provider.StartNewProvider(ctx, connector.GetName(), streamer, log)
}
