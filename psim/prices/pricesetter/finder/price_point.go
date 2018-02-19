package finder

import (
	"gitlab.com/swarmfund/psim/psim/prices/pricesetter/provider"
)

type providerPricePoint struct {
	ProviderID string
	provider.PricePoint
}
