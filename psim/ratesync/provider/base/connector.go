package base

import (
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"io/ioutil"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
	"gitlab.com/distributed_lab/logan/v3/fields"
)

type PricesResponse interface {
	ToPrices() ([]provider.PricePoint, error)
	fields.Provider
}

type Connector struct {
	Name string
	Endpoints map[string]string
}

// GetName retrieves the name of connector
func (c *Connector) GetName() string {
	return c.Name
}

func (c *Connector) GetPrices(baseAsset, quoteAsset string, pricesResponse PricesResponse) ([]provider.PricePoint, error) {
	assetPair := baseAsset + "/" + quoteAsset
	endpoint, ok := c.Endpoints[assetPair]
	if !ok {
		return nil, errors.From(errors.New("Unknown asset pair"), logan.F{
			"requested_asset_pair": assetPair,
		})
	}

	fieldz := logan.F{
		"endpoint": endpoint,
	}

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to perform get request", fieldz)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.From(errors.New("Unexpected status code."), fieldz.Merge(logan.F{
			"status_code": response.StatusCode,
		}))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read data from response body")
	}

	err = json.Unmarshal(body, pricesResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response body into pricesResponse", fieldz.Merge(logan.F{
			"body": string(body),
		}))
	}

	fieldz["prices_response"] = pricesResponse

	result, err := pricesResponse.ToPrices()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve Prices from the unmarshalled pricesResponse", fieldz)
	}
	return result, nil
}
