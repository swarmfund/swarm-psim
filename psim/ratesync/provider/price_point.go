package provider

import (
	"time"
	"gitlab.com/tokend/go/amount"
)

// PricePoint stores info on price value and time it was
type PricePoint struct {
	Price int64
	Time time.Time
}

// GetPercentDeltaToMinPrice - returns delta in percent compared to min price of two points
// returns false if overflows
// panics if min price is 0
func (p PricePoint) GetPercentDeltaToMinPrice(other PricePoint) (int64, bool) {
	minPrice := p.Price
	maxPrice := other.Price
	if other.Price < minPrice {
		minPrice = other.Price
		maxPrice = p.Price
	}

	return amount.BigDivide(maxPrice - minPrice, amount.One, minPrice, amount.ROUND_UP)
}
