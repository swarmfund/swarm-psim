package gdaxProvider

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

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

// StartNewGdaxProvider creates new gdaxProvider and runs it concurrently
func StartNewGdaxProvider(ctx context.Context, log *logan.Entry, baseAsset, quoteAsset string) <-chan price.PricePoint {
	provider := gdaxProvider{
		logger:        log.WithField("connector", "gdax"),
		baseAsset:     baseAsset,
		quoteAsset:    quoteAsset,
		pricesChannel: make(chan price.PricePoint),
	}

	go provider.runOnce(ctx)

	return provider.pricesChannel
}

func (p *gdaxProvider) runOnce(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			p.logger.WithError(errors.FromPanic(r)).Error("connector panicked")
			return
		}
	}()
	defer close(p.pricesChannel)

	assetPair := p.baseAsset + "-" + p.quoteAsset
	if !contains(assetPairs, assetPair) {
		p.logger.Error("unknown asset pair: ", assetPair)
		return
	}

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial(url, nil)
	if err != nil {
		p.logger.WithError(err).Error("failed to create a client connection")
		return
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
		p.logger.WithError(err).Error("failed to encode request message")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var jp jsonAssetPrice
			if err := wsConn.ReadJSON(&jp); err != nil {
				p.logger.WithError(err).Error("failed to decode response message")
				return
			}

			if jp.LastUpdated.Unix() <= 0 {
				continue
			}

			jps := jsonPrices{jp}
			prices, err := jps.Prices()
			if err != nil {
				p.logger.WithError(err).Error("failed to unmarshal prices")
				return
			}

			for _, item := range prices {
				p.pricesChannel <- item
			}
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
