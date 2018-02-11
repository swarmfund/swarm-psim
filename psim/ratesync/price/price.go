package price

import "time"

// PricePoint is a single price response
type PricePoint struct {
	Price int64
	Time  time.Time
}

// Prices is an array of PricePoint
type Prices []PricePoint
