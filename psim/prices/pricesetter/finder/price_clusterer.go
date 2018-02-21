package finder

type priceClusterizerImpl struct {
	points []providerPricePoint
}

func newPriceClusterizer(points []providerPricePoint) *priceClusterizerImpl {
	return &priceClusterizerImpl{
		points: points,
	}
}

type candidatePoint struct {
	providerPricePoint
	distance int64
}

// GetClusterForPoint for each Provider finds the Point among p.points, which is
// the closet to the provided point.
func (p *priceClusterizerImpl) GetClusterForPoint(point providerPricePoint) []providerPricePoint {
	providerToCandidate := map[string]candidatePoint{}

	for i := range p.points {
		// no need to process points of provider for which cluster is requested
		if p.points[i].ProviderID == point.ProviderID {
			continue
		}

		candidate := candidatePoint{
			providerPricePoint: p.points[i],
			distance:           calcDistance(point, p.points[i]),
		}

		bestPoint, ok := providerToCandidate[candidate.ProviderID]
		if !ok {
			providerToCandidate[candidate.ProviderID] = candidate
			continue
		}

		if bestPoint.distance > candidate.distance {
			// candidate if better than bestPoint - found new best Point
			providerToCandidate[candidate.ProviderID] = candidate
		}
	}

	result := []providerPricePoint{
		point,
	}

	for key := range providerToCandidate {
		result = append(result, providerToCandidate[key].providerPricePoint)
	}

	return result
}

func calcDistance(from, to providerPricePoint) int64 {
	if from.ProviderID == to.ProviderID {
		panic("Unexpected state: trying to calculate distance for points from same provider")
	}

	delta := from.Price - to.Price
	if delta < 0 {
		return -delta
	}

	return delta
}
