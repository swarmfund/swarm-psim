package bitstamp

import (
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider/base"
)

// Connector is a connector for bitfinex.com
type Connector struct {
	*base.Connector
}

// New is a constructor for Connector
func New() *Connector {
	return &Connector{
		Connector: &base.Connector{
			Name: "bitstamp",
			Endpoints: map[string]string{
				"BTC/USD": "https://www.bitstamp.net/api/v2/ticker/btcusd",
				"ETH/USD": "https://www.bitstamp.net/api/v2/ticker/ethusd",
			},
		},
	}
}

// GetPrices retrieves prices from external service and returns structured prices
func (c *Connector) GetPrices(baseAsset, quoteAsset string) ([]provider.PricePoint, error) {
	return c.Connector.GetPrices(baseAsset, quoteAsset, &jsonAssetPrice{})
}
