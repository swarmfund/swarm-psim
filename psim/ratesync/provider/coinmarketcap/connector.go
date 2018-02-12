package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/ratesync/price"
	"io/ioutil"
	"net/http"
	"time"
)

var assetPairs = map[string]string{
	"BTC/USD": "https://api.coinmarketcap.com/v1/ticker/bitcoin/?convert=USD",
	"ETH/USD": "https://api.coinmarketcap.com/v1/ticker/ethereum/?convert=USD",
	"SWM/USD": "https://api.coinmarketcap.com/v1/ticker/swarm-fund/?convert=USD",
}

// Connector is a connector for coinmarketcap.com
type Connector struct {
	Name string
}

// NewConnector is a constructor for Connector
func NewConnector() *Connector {
	return &Connector{
		Name: "coinmarketcap",
	}
}

// GetName retrieves the name of connector
func (c *Connector) GetName() string {
	return c.Name
}

// GetPrices retrieves prices from external service and returns structured prices
func (c *Connector) GetPrices(baseAsset, quoteAsset string) (price.Prices, error) {
	assetPair := baseAsset + "/" + quoteAsset
	if _, ok := assetPairs[assetPair]; !ok {
		return nil, fmt.Errorf("uknown asset pair: %v", assetPair)
	}

	request := assetPairs[assetPair]
	response, err := http.Get(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get price with status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read data from response body")
	}

	var jp jsonPrices
	err = json.Unmarshal(body, &jp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	result, err := jp.Prices()
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal prices")
	}
	return result, nil
}

// jsonAssetPrice is middle-layer structure for custom JSON unmarshalling
type jsonAssetPrice struct {
	PriceUsd    string `json:"price_usd"`
	LastUpdated int64  `json:"last_updated,string"`
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
			Time:  time.Unix(jp.LastUpdated, 0).UTC(),
		})
	}
	return result, nil
}
