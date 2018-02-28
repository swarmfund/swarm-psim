package providers

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/prices/providers/bitfinex"
	"gitlab.com/swarmfund/psim/psim/prices/providers/bitstamp"
	"gitlab.com/swarmfund/psim/psim/prices/providers/coinmarketcap"
	"gitlab.com/swarmfund/psim/psim/prices/providers/gdax"
)

func StartSpecificProvider(ctx context.Context, log *logan.Entry, baseAsset, quoteAsset, providerName string, providerPeriod time.Duration) (*Provider, error) {
	switch providerName {
	case bitfinex.Name:
		return priceProviderFromConnector(ctx, log, baseAsset, quoteAsset, bitfinex.New(), providerPeriod), nil
	case bitstamp.Name:
		return priceProviderFromConnector(ctx, log, baseAsset, quoteAsset, bitstamp.New(), providerPeriod), nil
	case coinmarketcap.Name:
		return priceProviderFromConnector(ctx, log, baseAsset, quoteAsset, coinmarketcap.New(), providerPeriod), nil
	case gdax.Name:
		// Gdax exchange provides Prices over socket, so Connector which does htt.Get over ticker is unsuitable here.
		pointsStream, err := gdax.StartNewPriceStreamer(ctx, log, baseAsset, quoteAsset)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create and start Gdax Streamer")
		}

		return StartNewProvider(ctx, providerName, pointsStream, log), nil
	default:
		return nil, errors.From(errors.New("Unexpected PriceProvider name"), logan.F{
			"provider_name": providerName,
		})
	}
}

func priceProviderFromConnector(
	ctx context.Context,
	log *logan.Entry,
	baseAsset string,
	quoteAsset string,
	connector Connector,
	period time.Duration) *Provider {

	pointsStream := StartNewPriceStreamer(ctx, log, baseAsset, quoteAsset, connector, period)
	return StartNewProvider(ctx, connector.GetName(), pointsStream, log)
}
