package gdax

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	ws "github.com/gorilla/websocket"
	"github.com/preichenberger/go-gdax"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/prices/types"
)

const (
	Name = "gdax"
)

var assetPairs = map[string]struct{}{
	"BTC-USD": {},
	"ETH-USD": {},
}

type streamer struct {
	logger        *logan.Entry
	assetPair     string
	pricesChannel chan types.PricePoint
}

// StartNewPriceStreamer creates new gdaxProvider and runs it safely and concurrently
func StartNewPriceStreamer(ctx context.Context, log *logan.Entry, baseAsset, quoteAsset string) (<-chan types.PricePoint, error) {
	assetPair := baseAsset + "-" + quoteAsset
	_, ok := assetPairs[assetPair]
	if !ok {
		return nil, errors.From(errors.New("Provided asset pair is not supported."), logan.F{
			"asset_pair":            assetPair,
			"supported_asset_pairs": assetPairs,
		})
	}

	p := streamer{
		logger:        log.WithField("prices_streamer", Name),
		assetPair:     assetPair,
		pricesChannel: make(chan types.PricePoint, 10),
	}

	// if runOnce returned - we have been disconnected from the provider, so it's better to wait before trying to connect again
	go app.RunOverIncrementalTimer(ctx, p.logger, Name, p.runOnce, time.Second*5, time.Minute)

	return p.pricesChannel, nil
}

func (p *streamer) runOnce(ctx context.Context) error {
	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		return errors.Wrap(err, "failed to create a client connection")
	}

	subscribe := gdax.Message{
		Type: "subscribe",
		Channels: []gdax.MessageChannel{
			{
				Name: "ticker",
				ProductIds: []string{
					p.assetPair,
				},
			},
		},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		return errors.Wrap(err, "failed to encode request message")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := p.read(ctx, wsConn)
			if err != nil {
				return errors.Wrap(err, "failed to read from connection")
			}
		}
	}

	return nil
}

func (p *streamer) read(ctx context.Context, wsConn *ws.Conn) error {
	var jp jsonAssetPrice
	if err := wsConn.ReadJSON(&jp); err != nil {
		return errors.Wrap(err, "failed to decode response message")
	}

	if jp.LastUpdated.Unix() <= 0 {
		return nil
	}

	jps := jsonPrices{jp}
	prices, err := jps.Prices()
	if err != nil {
		return errors.Wrap(err, "failed  to unmarshal prices")
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
