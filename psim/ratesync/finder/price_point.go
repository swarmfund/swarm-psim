package finder

import (
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
)

type providerPricePoint struct {
	ProviderID string
	provider.PricePoint
}
