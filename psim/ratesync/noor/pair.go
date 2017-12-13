package noor

const (
	ounce = 31.1034768
)

type Pair struct {
	Quote  string
	Code   string
	Symbol string
	Weight float64
}

func (p *Pair) PhysicalPrice(price float64) (int64) {
	return int64(price / ounce * p.Weight * 10000)
}
