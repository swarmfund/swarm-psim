package finder

import "gitlab.com/swarmfund/psim/psim/prices/pricesetter/providers"

type priceClusterizerImpl struct {
	providersPoints []providerPricePoint
}

func newPriceClusterizer(points []providerPricePoint) *priceClusterizerImpl {
	return &priceClusterizerImpl{
		providersPoints: points,
	}
}

type candidatePoint struct {
	providerPricePoint
	distance int64
}

// GetClusterForPoint for each Provider finds the Point among p.points, which is
// the closet to the provided point.
func (p *priceClusterizerImpl) GetClusterForPoint(point providers.PricePoint) []providerPricePoint {
	providerToCandidate := map[string]candidatePoint{}

	for i := range p.providersPoints {
		candidate := candidatePoint{
			providerPricePoint: p.providersPoints[i],
			distance:           calcDistance(point, p.providersPoints[i].PricePoint),
		}

		bestPoint, ok := providerToCandidate[candidate.ProviderID]
		if !ok {
			// Still no best Point for this Provider
			providerToCandidate[candidate.ProviderID] = candidate
			continue
		}

		if bestPoint.distance > candidate.distance {
			// candidate if better than bestPoint for the Provider - found new best Point
			providerToCandidate[candidate.ProviderID] = candidate
		}
	}

	var result []providerPricePoint

	for key := range providerToCandidate {
		result = append(result, providerToCandidate[key].providerPricePoint)
	}

	return result
}

func calcDistance(p1, p2 providers.PricePoint) int64 {
	delta := p1.Price - p2.Price
	if delta < 0 {
		return -delta
	}

	return delta
}
