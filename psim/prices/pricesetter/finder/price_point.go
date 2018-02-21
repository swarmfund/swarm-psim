package finder

import (
	"gitlab.com/swarmfund/psim/psim/prices/pricesetter/provider"
)

type providerPricePoint struct {
	ProviderID string
	provider.PricePoint
}

func (p providerPricePoint) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"provider_id": p.ProviderID,
		"point":       p.PricePoint,
	}
}
