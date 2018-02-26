package finder

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"sort"
	"time"
)

//go:generate mockery -case underscore -testonly -inpkg -name ratesProvider
type RatesProvider interface {
	// GetName - returns name of the provider
	GetName() string
	// GetPoints - returns points available
	GetPoints() []provider.PricePoint
	// RemoveDeprecatedPoints - removes points which were created before minAllowedTime
	RemoveDeprecatedPoints(minAllowedTime time.Time)
}

//go:generate mockery -case underscore -testonly -inpkg -name priceClusterer
type priceClusterer interface {
	// GetClusterForPoint - returns cluster of nearest points for specified point
	GetClusterForPoint(point pricePoint) []pricePoint
}

type priceClustererProvider func(points []pricePoint) priceClusterer

type priceFinder struct {
	log            *logan.Entry
	ratesProviders []RatesProvider
	priceClustererProvider

	maxPercentPriceDelta        int64
	minNumberOfProvidersToAgree int
}

func NewPriceFinder(log *logan.Entry, ratesProviders []RatesProvider, maxPercentPriceDelta int64,
	minNumberOfProvidersToAgree int) (*priceFinder, error) {
	if !isPercentValid(maxPercentPriceDelta) {
		return nil, errors.New("maxPercentPriceDelta must be in range [0, 100*ONE]")
	}

	if minNumberOfProvidersToAgree <= 0 {
		return nil, errors.New("minNumberOfProvidersToAgree must be > 0")
	}

	if len(ratesProviders) == 0 {
		return nil, errors.New("Unexpected number of rate providers")
	}

	result := &priceFinder{
		ratesProviders:              ratesProviders,
		maxPercentPriceDelta:        maxPercentPriceDelta,
		minNumberOfProvidersToAgree: minNumberOfProvidersToAgree,
		priceClustererProvider: priceClustererProvider(func(points []pricePoint) priceClusterer {
			return newPriceClusterer(points)
		}),
		log: log,
	}

	return result, nil
}

func isPercentValid(percent int64) bool {
	return 0 <= percent && percent <= 100*amount.One
}

// TryFind - tries to find most recent price point which priceDelta is <= maxPercentPriceDelta
// for percent of providers >= minPercentOfNodeToParticipate
func (p *priceFinder) TryFind() (*provider.PricePoint, error) {
	allPoints, err := p.getAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get points available")
	}

	sort.Sort(sortablePricePointsByTime(allPoints))
	clusterer := p.priceClustererProvider(allPoints)

	for i := range allPoints {
		cluster := clusterer.GetClusterForPoint(allPoints[i])
		candidate := median(cluster)
		if p.isMeetsReq(candidate, cluster) {
			p.log.WithField("point", candidate).WithField("cluster", cluster).Info("Found point meeting restrictions")
			return &candidate.PricePoint, nil
		}
	}

	p.log.Info("Failed to find point meeting restrictions")
	return nil, nil
}

// RemoveDeprecatedPoints - removes points which were created before minAllowedTime
func (p *priceFinder) RemoveDeprecatedPoints(minAllowedTime time.Time) {
	for i := range p.ratesProviders {
		p.ratesProviders[i].RemoveDeprecatedPoints(minAllowedTime)
	}
}

func median(data []pricePoint) pricePoint {
	sort.Sort(sortablePricePointsByPrice(data))
	// we can not select mean from two median points as point must exist
	return data[len(data)/2]
}

func (p *priceFinder) getAll() ([]pricePoint, error) {
	var result []pricePoint
	for _, ratesProvider := range p.ratesProviders {
		pricePoints := ratesProvider.GetPoints()

		p.log.WithField("provider", ratesProvider.GetName()).WithField("points", pricePoints).Info("Getting points to find consensus on new one")
		for i := range pricePoints {
			result = append(result, pricePoint{
				ProviderID: ratesProvider.GetName(),
				PricePoint: pricePoints[i],
			})
		}

	}

	return result, nil
}

func (p *priceFinder) isMeetsReq(candidate pricePoint, cluster []pricePoint) bool {
	if p.minNumberOfProvidersToAgree > len(cluster) {
		p.log.WithField("min number of providers", p.minNumberOfProvidersToAgree).WithField("cluster size", len(cluster)).
			Debug("Cluster too small - skipping")
		return false
	}

	providersAgreed := 0
	for _, providerPricePoint := range cluster {
		if p.isMeetsPriceDeltaReq(candidate, providerPricePoint) {
			p.log.WithField("providerPricePoint", providerPricePoint).Debug("Agreed")
			providersAgreed++
		} else {
			p.log.WithField("providerPricePoint", providerPricePoint).Debug("Disagreed")
		}
	}

	return providersAgreed >= p.minNumberOfProvidersToAgree
}

func (p *priceFinder) isMeetsPriceDeltaReq(candidate, point pricePoint) bool {
	percentInDelta, isOverflow := candidate.GetPercentDeltaToMinPrice(point.PricePoint)
	if isOverflow {
		p.log.WithField("candidate_price", candidate.Price).
			WithField("point_price", point.Price).Warn("Overflow on price delta calculation")
		return false
	}

	return percentInDelta <= p.maxPercentPriceDelta
}
