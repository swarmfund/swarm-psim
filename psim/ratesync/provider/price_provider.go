package provider

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/ratesync/price"
	"time"
)

// Connector is an interface for retrieving asset prices from external services
type Connector interface {
	GetName() string
	GetPrices(baseAsset, quoteAsset string) (price.Prices, error)
}

// pricesProvider provides concurrent work of multiple connectors
type pricesProvider struct {
	logger        *logan.Entry
	baseAsset     string
	quoteAsset    string
	pricesChannel chan price.PricePoint
	exchange      Connector
	period        time.Duration
}

// StartNewPricesProvider creates new pricesProvider and runs it safely and concurrently
func StartNewPricesProvider(ctx context.Context, log *logan.Entry, baseAsset, quoteAsset string, exchange Connector, period time.Duration) <-chan price.PricePoint {
	priceProvider := pricesProvider{
		logger:        log.WithField("connector", exchange.GetName()),
		baseAsset:     baseAsset,
		quoteAsset:    quoteAsset,
		pricesChannel: make(chan price.PricePoint),
		exchange:      exchange,
		period:        period,
	}

	go priceProvider.runSafely(ctx)

	return priceProvider.pricesChannel
}

func (p *pricesProvider) runOnce(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrap(errors.WithStack(errors.FromPanic(r)), "connector panicked")
			return
		}
	}()

	var prices price.Prices
	prices, err = p.exchange.GetPrices(p.baseAsset, p.quoteAsset)
	if err != nil {
		return errors.Wrap(err, "failed to get prices")
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

func (p *pricesProvider) runSafely(ctx context.Context) {
	t := time.NewTicker(p.period)
	defer close(p.pricesChannel)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := p.runOnce(ctx)
			if err != nil {
				p.logger.WithError(err).Error("prices provider returned error")
			}
		}
	}
}
