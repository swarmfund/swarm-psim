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
type ratesProvider interface {
	// GetName - returns name of the provider
	GetName() string
	// GetPoints - returns points available
	GetPoints() ([]provider.PricePoint, error)
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
	ratesProviders []ratesProvider
	priceClustererProvider

	maxPercentPriceDelta         int64
	minPercentOfProvidersToAgree int64
}

func newPriceFinder(log *logan.Entry, ratesProviders []ratesProvider, maxPercentPriceDelta int64,
	minPercentOfProvidersToAgree int64, priceClustererProvider priceClustererProvider) (*priceFinder, error) {
	if !isPercentValid(maxPercentPriceDelta) {
		return nil, errors.New("maxPercentPriceDelta must be in range [0, 100*ONE]")
	}

	if !isPercentValid(minPercentOfProvidersToAgree) {
		return nil, errors.New("minPercentOfProvidersToAgree must be in range [0, 100*ONE]")
	}

	if len(ratesProviders) == 0 {
		return nil, errors.New("Unexpected number of rate providers")
	}

	result := &priceFinder{
		ratesProviders:               ratesProviders,
		maxPercentPriceDelta:         maxPercentPriceDelta,
		minPercentOfProvidersToAgree: minPercentOfProvidersToAgree,
		priceClustererProvider:       priceClustererProvider,
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
			return &candidate.PricePoint, nil
		}
	}

	return nil, nil
}

func median(data []pricePoint) pricePoint {
	sort.Sort(sortablePricePointsByPrice(data))
	// we can not select mean from two median points as point must exist
	return data[len(data)/2]
}

func (p *priceFinder) getAll() ([]pricePoint, error) {
	var result []pricePoint
	for _, ratesProvider := range p.ratesProviders {
		pricePoints, err := ratesProvider.GetPoints()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get points", logan.F{
				"rate_provider": ratesProvider.GetName(),
			})
		}

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
	providersAgreed := int64(0)
	for _, providerPricePoints := range cluster {
		if p.isMeetsPriceDeltaReq(candidate, providerPricePoints) {
			providersAgreed++
		}
	}

	totalNumberOfProviders := int64(len(p.ratesProviders))
	percentOfProvidersAgreed, isOverflow := amount.BigDivide(providersAgreed, 100*amount.One, totalNumberOfProviders, amount.ROUND_DOWN)
	if isOverflow {
		p.log.WithField("agreed", providersAgreed).WithField("total", totalNumberOfProviders).
			Warn("Overflow on percentOfProvidersAgreed calculation")
		return false
	}

	return percentOfProvidersAgreed >= p.minPercentOfProvidersToAgree
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
