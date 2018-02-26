package providers

import (
	"time"
	"gitlab.com/swarmfund/go/amount"
)

// PricePoint stores info on Price value and time it was
type PricePoint struct {
	Price int64
	Time time.Time
}

// GetPercentDeltaToMinPrice - returns delta in percent compared to min price of two points
//
// returns false if overflows
// panics if min price is 0
func (p PricePoint) GetPercentDeltaToMinPrice(other PricePoint) (int64, bool) {
	minPrice := p.Price
	maxPrice := other.Price
	if other.Price < minPrice {
		minPrice = other.Price
		maxPrice = p.Price
	}

	// delta * amount.One / minPrice
	return amount.BigDivide(maxPrice - minPrice, amount.One, minPrice, amount.ROUND_UP)
}

func (p PricePoint) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"price": p.Price,
		"time":   p.Time.String(),
	}
}
