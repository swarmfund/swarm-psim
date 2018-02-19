package bitfinex

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
)

// jsonAssetPrice is middle-layer structure for custom JSON unmarshalling
type jsonAssetPrice struct {
	PriceUsd    string  `json:"last_price"`
	LastUpdated float64 `json:"timestamp,string"`
}

func (jp jsonAssetPrice) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"price_usd":    jp.PriceUsd,
		"last_updated": jp.LastUpdated,
	}
}

// ToPrices returns unmarshaled array of PricePoint with appropriate representation of price and time
func (jp *jsonAssetPrice) ToPrices() ([]provider.PricePoint, error) {
	price, err := amount.Parse(jp.PriceUsd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse amount", logan.F{
			"raw_amount": jp.PriceUsd,
		})
	}

	return []provider.PricePoint{
		{
			Price: price,
			Time:  time.Unix(int64(jp.LastUpdated), 0).UTC(),
		},
	}, nil
}
