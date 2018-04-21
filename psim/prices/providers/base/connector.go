package base

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3/fields"
	"gitlab.com/swarmfund/psim/psim/prices/types"
)

type PricesResponse interface {
	ToPrices() ([]types.PricePoint, error)
	fields.Provider
}

type Connector struct {
	Name string
	// TODO Consider using some new specific AssetPair type as a key.
	Endpoints map[string]string
}

// GetName retrieves the name of connector
func (c *Connector) GetName() string {
	return c.Name
}

func (c *Connector) GetPrices(baseAsset, quoteAsset string, pricesResponse PricesResponse) ([]types.PricePoint, error) {
	assetPair := baseAsset + "/" + quoteAsset
	endpoint, ok := c.Endpoints[assetPair]
	if !ok {
		return nil, errors.From(errors.New("Unknown asset pair."), logan.F{
			"requested_asset_pair": assetPair,
		})
	}

	fields := logan.F{
		"endpoint": endpoint,
	}

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to perform get request", fields)
	}
	fields["status_code"] = response.StatusCode

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.From(errors.New("Unexpected status code."), fields)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read data from response body", fields)
	}

	err = json.Unmarshal(body, pricesResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response body into pricesResponse", fields.Merge(logan.F{
			"raw_response_body": string(body),
		}))
	}

	fields["prices_response"] = pricesResponse

	result, err := pricesResponse.ToPrices()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve Prices from the unmarshalled pricesResponse", fields)
	}

	return result, nil
}
