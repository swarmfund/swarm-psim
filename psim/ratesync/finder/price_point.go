package finder

import (
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
)

type pricePoint struct {
	ProviderID string
	provider.PricePoint
}
