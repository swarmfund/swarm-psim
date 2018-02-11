package bitstamp

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
	"BTC/USD": "https://www.bitstamp.net/api/v2/ticker/btcusd",
	"ETH/USD": "https://www.bitstamp.net/api/v2/ticker/ethusd",
}

// BitstampConnector is a connector for bitstamp.net
type BitstampConnector struct {
	Name string
}

// NewBitstampConnector is a constructor for BitstampConnector
func NewBitstampConnector() *BitstampConnector {
	return &BitstampConnector{
		Name: "bitstamp",
	}
}

// GetName retrieves the name of connector
func (c *BitstampConnector) GetName() string {
	return c.Name
}

// GetPrices retrieves prices from external service and returns structured prices
func (c *BitstampConnector) GetPrices(baseAsset, quoteAsset string) (price.Prices, error) {
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

	var jp jsonAssetPrice
	err = json.Unmarshal(body, &jp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}
	jps := jsonPrices{jp}

	result, err := jps.Prices()
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal prices")
	}
	return result, nil
}

// jsonAssetPrice is middle-layer structure for custom JSON unmarshalling
type jsonAssetPrice struct {
	PriceUsd    string `json:"last"`
	LastUpdated int64  `json:"timestamp,string"`
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
			Time:  time.Unix(jp.LastUpdated, 0),
		})
	}
	return result, nil
}
