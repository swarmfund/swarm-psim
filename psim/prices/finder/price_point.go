package finder

import (
	"gitlab.com/swarmfund/psim/psim/prices/types"
)

type providerPricePoint struct {
	ProviderID string
	types.PricePoint
}

func (p providerPricePoint) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"provider_id": p.ProviderID,
		"point":       p.PricePoint,
	}
}
