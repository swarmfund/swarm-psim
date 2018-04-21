package providers

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/prices/types"
	"gitlab.com/distributed_lab/running"
)

// Connector is an interface for retrieving asset prices from external services
type Connector interface {
	GetName() string
	GetPrices(baseAsset, quoteAsset string) ([]types.PricePoint, error)
}

// Streamer obtains types.PricePoints with the GetPrices method of the Connector from time to time
// and streams these types.PricePoints into the pricesChannel.
//
// Streamer is a common Streamer of Prices for a PriceProvider.
type streamer struct {
	log           *logan.Entry
	baseAsset     string
	quoteAsset    string
	pricesChannel chan types.PricePoint
	exchange      Connector
	period        time.Duration
}

// StartNewPriceStreamer creates new Streamer and runs it safely and concurrently.
// StartNewPriceStreamer is *not* a blocking function.
func StartNewPriceStreamer(
	ctx context.Context,
	log *logan.Entry,
	baseAsset,
	quoteAsset string,
	exchange Connector,
	period time.Duration) <-chan types.PricePoint {

	streamer := streamer{
		log: log.WithFields(logan.F{
			"price_streamer": exchange.GetName(),
			"base_asset":     baseAsset,
			"quote_asset":    quoteAsset,
		}),
		baseAsset:     baseAsset,
		quoteAsset:    quoteAsset,
		pricesChannel: make(chan types.PricePoint, 10),
		exchange:      exchange,
		period:        period,
	}

	streamer.log.Debug("Starting new PriceStreamer.")
	go running.WithBackOff(ctx, streamer.log, exchange.GetName(), streamer.runOnce, period, 10*time.Second, time.Hour)

	return streamer.pricesChannel
}

func (p *streamer) runOnce(ctx context.Context) error {
	prices, err := p.exchange.GetPrices(p.baseAsset, p.quoteAsset)
	if err != nil {
		return errors.Wrap(err, "Failed to get Prices from ExchangeConnector")
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
