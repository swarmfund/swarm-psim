package finder

type priceClustererImpl struct {
	points []pricePoint
}

func newPriceClusterer(points []pricePoint) *priceClustererImpl {
	return &priceClustererImpl{
		points: points,
	}
}

type candidatePoint struct {
	pricePoint
	distance int64
}

// GetClusterForPoint - returns cluster of nearest points for specified point
func (p *priceClustererImpl) GetClusterForPoint(point pricePoint) []pricePoint {
	candidates := map[string]candidatePoint{}

	for i := range p.points {
		// no need to process points of provider for which cluster is requested
		if p.points[i].ProviderID == point.ProviderID {
			continue
		}


		candidate := candidatePoint{
			pricePoint: p.points[i],
			distance: calcDistance(point, p.points[i]),
		}

		bestPoint, ok := candidates[candidate.ProviderID]
		if !ok {
			candidates[candidate.ProviderID] = candidate
			continue
		}

		if bestPoint.distance > candidate.distance {
			candidates[candidate.ProviderID] = candidate
		}
	}

	result := []pricePoint{
		point,
	}

	for key := range candidates {
		result = append(result, candidates[key].pricePoint)
	}

	return result
}

func calcDistance(from, to pricePoint) (int64) {
	if from.ProviderID == to.ProviderID {
		panic("Unexpected state: trying to calculate distance for points from same provider")
	}

	delta := from.Price - to.Price
	if delta < 0 {
		return -delta
	}

	return delta
}

