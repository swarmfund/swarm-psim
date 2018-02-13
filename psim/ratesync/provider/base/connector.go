package base

import (
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"io/ioutil"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
)

type Prices interface {
	ToPrices() ([]provider.PricePoint, error)
}

type Connector struct {
	Name string
	Endpoints map[string]string
}

// GetName retrieves the name of connector
func (c *Connector) GetName() string {
	return c.Name
}

func (c *Connector) GetPrices(baseAsset, quoteAsset string, prices Prices) ([]provider.PricePoint, error) {
	assetPair := baseAsset + "/" + quoteAsset
	endpoint, ok := c.Endpoints[assetPair]
	if !ok {
		return nil, errors.From(errors.New("Unknown asset pair"), logan.F{
			"requested_asset_pair": assetPair,
		})
	}

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform get request")
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.From(errors.New("Unexpected status code"), logan.F{
			"status_code": response.StatusCode,
		})
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read data from response body")
	}

	err = json.Unmarshal(body, prices)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body", logan.F{
			"body": string(body),
		})
	}

	result, err := prices.ToPrices()
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal prices")
	}
	return result, nil
}
