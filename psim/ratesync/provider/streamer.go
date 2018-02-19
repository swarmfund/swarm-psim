package provider

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"time"
	"gitlab.com/swarmfund/psim/psim/app"
)

// Connector is an interface for retrieving asset prices from external services
type Connector interface {
	GetName() string
	GetPrices(baseAsset, quoteAsset string) ([]PricePoint, error)
}

// Streamer obtains PricePoints with the GetPrices method of the Connector from time to time
// and streams these PricePoints into the pricesChannel.
//
// Streamer is a common Streamer of Prices for a PriceProvider.
type streamer struct {
	logger        *logan.Entry
	baseAsset     string
	quoteAsset    string
	pricesChannel chan PricePoint
	exchange      Connector
	period        time.Duration
}

// StartNewPriceStreamer creates new Streamer and runs it safely and concurrently
func StartNewPriceStreamer(
	ctx context.Context,
	log *logan.Entry,
	baseAsset,
	quoteAsset string,
	exchange Connector,
	period time.Duration) <-chan PricePoint {

	streamer := streamer{
		logger:        log.WithField("price_streamer", exchange.GetName()),
		baseAsset:     baseAsset,
		quoteAsset:    quoteAsset,
		pricesChannel: make(chan PricePoint, 10),
		exchange:      exchange,
		period:        period,
	}

	streamer.logger.Debug("Starting new Streamer.")
	go app.RunOverIncrementalTimer(ctx, streamer.logger, exchange.GetName(), streamer.runOnce, period, time.Minute)

	return streamer.pricesChannel
}

func (p *streamer) runOnce(ctx context.Context) (err error) {
	var prices []PricePoint

	prices, err = p.exchange.GetPrices(p.baseAsset, p.quoteAsset)
	if err != nil {
		return errors.Wrap(err, "Failed to get Prices from Exchange Connector")
	}

	for _, item := range prices {
		select {
		case p.pricesChannel <- item:
			continue
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}
