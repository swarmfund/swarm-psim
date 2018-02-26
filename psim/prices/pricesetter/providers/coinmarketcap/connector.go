package coinmarketcap

import (
	"gitlab.com/swarmfund/psim/psim/prices/pricesetter/providers"
	"gitlab.com/swarmfund/psim/psim/prices/pricesetter/providers/base"
)

const (
	Name = "coinmarketcap"
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
				"BTC/USD": "https://api.coinmarketcap.com/v1/ticker/bitcoin/?convert=USD",
				"ETH/USD": "https://api.coinmarketcap.com/v1/ticker/ethereum/?convert=USD",
				"SWM/USD": "https://api.coinmarketcap.com/v1/ticker/swarm-fund/?convert=USD",
			},
		},
	}
}

// GetPrices retrieves prices from external service and returns structured prices
func (c *Connector) GetPrices(baseAsset, quoteAsset string) ([]providers.PricePoint, error) {
	return c.Connector.GetPrices(baseAsset, quoteAsset, &jsonPrices{})
}
