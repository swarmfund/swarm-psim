package gdax

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
)

// jsonAssetPrice is middle-layer structure for custom JSON unmarshalling
type jsonAssetPrice struct {
	Price       string    `json:"price"`
	LastUpdated time.Time `json:"time,string"`
}

// jsonPrices is an array of jsonAssetPrice
type jsonPrices []jsonAssetPrice

// Prices returns unmarshaled array of PricePoint with appropriate representation of price and time
func (jps jsonPrices) Prices() ([]providers.PricePoint, error) {
	var result []providers.PricePoint
	for _, jp := range jps {
		p, err := amount.Parse(jp.Price)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse amount")
		}

		result = append(result, providers.PricePoint{
			Price: p,
			Time:  jp.LastUpdated,
		})
	}
	return result, nil
}
