package coinmarketcap

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/prices/pricesetter/provider"
	"time"
)

// jsonAssetPrice is middle-layer structure for custom JSON unmarshalling
type jsonAssetPrice struct {
	PriceUsd    string `json:"price_usd"`
	LastUpdated int64  `json:"last_updated,string"`
}

func (jp jsonAssetPrice) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"price_usd":    jp.PriceUsd,
		"last_updated": jp.LastUpdated,
	}
}

// JsonPrices is an array of jsonAssetPrice
type jsonPrices []jsonAssetPrice

// ToPrices returns unmarshaled slice of PricePoint with appropriate representation of Price and time
func (jps jsonPrices) ToPrices() ([]provider.PricePoint, error) {
	var result []provider.PricePoint
	for _, jp := range jps {
		p, err := amount.Parse(jp.PriceUsd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse amount")
		}

		result = append(result, provider.PricePoint{
			Price: p,
			Time:  time.Unix(jp.LastUpdated, 0).UTC(),
		})
	}
	return result, nil
}

func (jps jsonPrices) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"prices": jps,
	}
}
