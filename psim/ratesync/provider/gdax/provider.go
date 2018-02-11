package gdaxCondnector

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"fmt"
	ws "github.com/gorilla/websocket"
	"github.com/preichenberger/go-gdax"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/ratesync/price"
)

const url = "wss://ws-feed.gdax.com"

var assetPairs = []string{
	"BTC-USD",
	"ETH-USD",
}

type gdaxProvider struct {
	logger        *logan.Entry
	baseAsset     string
	quoteAsset    string
	pricesChannel chan price.PricePoint
}

// StartNewGdaxProvider creates new gdaxProvider and runs it safely and concurrently
func StartNewGdaxProvider(ctx context.Context, log *logan.Entry, baseAsset, quoteAsset string) <-chan price.PricePoint {
	gdaxProvider := gdaxProvider{
		logger:        log.WithField("connector", "gdax"),
		baseAsset:     baseAsset,
		quoteAsset:    quoteAsset,
		pricesChannel: make(chan price.PricePoint),
	}

	go gdaxProvider.runSafely(ctx)

	return gdaxProvider.pricesChannel
}

func (p *gdaxProvider) runOnce(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrap(errors.WithStack(errors.FromPanic(r)), "connector panicked")
			return
		}
	}()
	defer close(p.pricesChannel)

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

func (p *gdaxProvider) GetPrices(baseAsset, quoteAsset string) (<-chan price.PricePoint, error) {
	assetPair := baseAsset + "-" + quoteAsset
	if !contains(assetPairs, assetPair) {
		return nil, fmt.Errorf("uknown asset pair: %v", assetPair)
	}

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a client connection")
	}

	subscribe := gdax.Message{
		Type: "subscribe",
		Channels: []gdax.MessageChannel{
			{
				Name: "ticker",
				ProductIds: []string{
					assetPair,
				},
			},
		},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		return nil, errors.Wrap(err, "failed to encode request message")
	}

	for {
		var jp jsonAssetPrice
		if err := wsConn.ReadJSON(&jp); err != nil {
			return nil, errors.Wrap(err, "failed to decode response message")
		}

		if jp.LastUpdated.Unix() <= 0 {
			continue
		}

		jps := jsonPrices{jp}
		prices, err := jps.Prices()
		if err != nil {
			return nil, errors.Wrap(err, "failed  to unmarshal prices")
		}

		for _, item := range prices {
			p.PricesChannel <- item
		}
	}
}

// jsonAssetPrice is middle-layer structure for custom JSON unmarshalling
type jsonAssetPrice struct {
	PriceUsd    string    `json:"price"`
	LastUpdated time.Time `json:"time,string"`
}

// jsonPrices is an array of jsonAssetPrice
type jsonPrices []jsonAssetPrice

// Prices returns unmarshaled array of PricePoint with appropriate representation of price and time
func (jps jsonPrices) Prices() (price.Prices, error) {
	var result price.Prices
	for _, jp := range jps {
		p, err := amount.Parse(jp.PriceUsd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse amount")
		}

		result = append(result, price.PricePoint{
			Price: p,
			Time:  jp.LastUpdated,
		})
	}
	return result, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
