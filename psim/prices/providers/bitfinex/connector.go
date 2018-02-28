package bitfinex

import (
	"gitlab.com/swarmfund/psim/psim/prices/providers/base"
	"gitlab.com/swarmfund/psim/psim/prices/types"
)

const (
	Name = "bitfinex"
)

// Connector is a connector for bitfinex.com
type Connector struct {
	*base.Connector
}

// New is a constructor for Connector
func New() *Connector {
	return &Connector{
		Connector: &base.Connector{
			Name: Name,
			Endpoints: map[string]string{
				"BTC/USD": "https://api.bitfinex.com/v1/pubticker/btcusd",
				"ETH/USD": "https://api.bitfinex.com/v1/pubticker/ethusd",
			},
		},
	}
}

// GetPrices retrieves prices from external service and returns structured prices
func (c *Connector) GetPrices(baseAsset, quoteAsset string) ([]types.PricePoint, error) {
	return c.Connector.GetPrices(baseAsset, quoteAsset, &jsonAssetPrice{})
}
