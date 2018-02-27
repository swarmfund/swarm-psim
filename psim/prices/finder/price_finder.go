package finder

import (
	"sort"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/prices/types"
)

//go:generate mockery -case underscore -testonly -inpkg -name PriceProvider
type PriceProvider interface {
	// GetName - returns name of the Provider
	GetName() string
	// GetPoints - returns Points available
	GetPoints() []types.PricePoint
	// RemoveOldPoints - removes points which were created before minAllowedTime
	RemoveOldPoints(minAllowedTime time.Time)
}

//go:generate mockery -case underscore -testonly -inpkg -name priceClusterizer
type priceClusterizer interface {
	// GetClusterForPoint - returns cluster of nearest Points for specified Point.
	// The provided point must *not* be included into the cluster.
	GetClusterForPoint(point types.PricePoint) []providerPricePoint
}

type priceClusterizerProvider func(points []providerPricePoint) priceClusterizer

type priceFinder struct {
	log            *logan.Entry
	priceProviders []PriceProvider
	priceClusterizerProvider

	//in range [0, 100*ONE]
	maxPercentPriceDelta        int64
	minNumberOfProvidersToAgree int
}

func NewPriceFinder(
	log *logan.Entry,
	priceProviders []PriceProvider,
	maxPercentPriceDelta int64,
	minNumberOfProvidersToAgree int) (*priceFinder, error) {

	if !isPercentValid(maxPercentPriceDelta) {
		return nil, errors.New("maxPercentPriceDelta must be in range [0, 100*ONE]")
	}

	if minNumberOfProvidersToAgree <= 0 {
		return nil, errors.New("minNumberOfProvidersToAgree must be > 0")
	}

	if len(priceProviders) == 0 {
		return nil, errors.New("Unexpected number of PriceProviders.")
	}

	result := &priceFinder{
		log:            log,
		priceProviders: priceProviders,
		priceClusterizerProvider: func(points []providerPricePoint) priceClusterizer {
			return newPriceClusterizer(points)
		},
		maxPercentPriceDelta:        maxPercentPriceDelta,
		minNumberOfProvidersToAgree: minNumberOfProvidersToAgree,
	}

	return result, nil
}

func isPercentValid(percent int64) bool {
	return percent >= 0 && percent <= 100*amount.One
}

// TryFind - tries to find most recent PricePoint which priceDelta is <= maxPercentPriceDelta
// for percent of providers >= minPercentOfNodeToParticipate
func (p *priceFinder) TryFind() (*types.PricePoint, error) {
	allPoints := p.getAllProvidersPoints()

	sort.Sort(sortablePricePointsByTime(allPoints))
	clusterizer := p.priceClusterizerProvider(allPoints)

	for i := range allPoints {
		cluster := clusterizer.GetClusterForPoint(allPoints[i].PricePoint)
		candidate := median(cluster)

		if p.meetsReq(candidate, cluster) {
			p.log.WithFields(logan.F{
				"point":   candidate,
				"cluster": cluster,
			}).Info("Found Point meeting restrictions.")
			return &candidate.PricePoint, nil
		}
	}

	p.log.Info("Hasn't found PricePoint meeting restrictions.")
	return nil, nil
}

// IsOK returns true if median of the Cluster built over all existing Points
// meets delta requirements with the provided point.
func (p *priceFinder) IsOK(point types.PricePoint) bool {
	allPoints := p.getAllProvidersPoints()
	clusterizer := p.priceClusterizerProvider(allPoints)

	cluster := clusterizer.GetClusterForPoint(point)
	clusterMedianPoint := median(cluster).PricePoint

	return p.meetsPriceDeltaReq(point, clusterMedianPoint)
}

// RemoveOldPoints - removes points which were created before minAllowedTime
func (p *priceFinder) RemoveOldPoints(minAllowedTime time.Time) {
	for i := range p.priceProviders {
		p.priceProviders[i].RemoveOldPoints(minAllowedTime)
	}
}

func (p *priceFinder) getAllProvidersPoints() []providerPricePoint {
	var result []providerPricePoint

	for _, priceProvider := range p.priceProviders {
		pricePoints := priceProvider.GetPoints()

		p.log.WithFields(logan.F{
			"provider": priceProvider.GetName(),
			"points":   pricePoints,
		}).Debug("Getting PricePoints to find consensus on new one.")

		for i := range pricePoints {
			result = append(result, providerPricePoint{
				ProviderID: priceProvider.GetName(),
				PricePoint: pricePoints[i],
			})
		}
	}

	return result
}

// Median returns one of the provided points.
func median(points []providerPricePoint) providerPricePoint {
	sort.Sort(sortablePricePointsByPrice(points))
	// we can not select mean from two median points as point must exist
	return points[len(points)/2]
}

// IsMeetsReq shows whether a candidate meets requirements.
// Returns true if cluster contains at least minNumberOfProvidersToAgree PricePoints, which agree with candidate.
func (p *priceFinder) meetsReq(candidate providerPricePoint, cluster []providerPricePoint) bool {
	if len(cluster) < p.minNumberOfProvidersToAgree {
		p.log.WithFields(logan.F{
			"min_number_of_providers": p.minNumberOfProvidersToAgree,
			"cluster_size":            len(cluster),
			"cluster":                 cluster,
		}).Debug("Cluster too small - skipping.")
		return false
	}

	providersAgreed := 0
	for _, providerPricePoint := range cluster {
		if p.meetsPriceDeltaReq(candidate.PricePoint, providerPricePoint.PricePoint) {
			// Provider agreed with the candidate - diff between the providerPricePoint and the candidate isn't too big.
			p.log.WithField("provider_price_point", providerPricePoint).Debug("Provider agreed.")
			providersAgreed++
		} else {
			p.log.WithField("provider_price_point", providerPricePoint).Debug("Provider disagreed.")
		}
	}

	return providersAgreed >= p.minNumberOfProvidersToAgree
}

func (p *priceFinder) meetsPriceDeltaReq(candidate, point types.PricePoint) bool {
	percentInDelta, isOverflow := candidate.GetPercentDeltaToMinPrice(point)
	if isOverflow {
		p.log.WithFields(logan.F{
			"candidate_price": candidate.Price,
			"point_price":     point.Price,
		}).Warn("Overflow on price delta calculation.")
		return false
	}

	return percentInDelta <= p.maxPercentPriceDelta
}
