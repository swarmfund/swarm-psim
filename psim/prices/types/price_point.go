package types

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
func GetPercentDeltaToMinPrice(price1, price2 int64) (int64, bool) {
	minPrice := price1
	maxPrice := price2
	if price2 < minPrice {
		minPrice = price2
		maxPrice = price1
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
