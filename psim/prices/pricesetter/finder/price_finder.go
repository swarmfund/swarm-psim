package finder

import (
	"sort"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/prices/pricesetter/provider"
)

// TODO Rename ratesProvider to pricesProvider
//go:generate mockery -case underscore -testonly -inpkg -name ratesProvider
type PriceProvider interface {
	// GetName - returns name of the Provider
	GetName() string
	// GetPoints - returns Points available
	GetPoints() []provider.PricePoint
	// RemoveOldPoints - removes points which were created before minAllowedTime
	RemoveOldPoints(minAllowedTime time.Time)
}

//go:generate mockery -case underscore -testonly -inpkg -name priceClusterer
type priceClusterer interface {
	// GetClusterForPoint - returns cluster of nearest points for specified point
	GetClusterForPoint(point providerPricePoint) []providerPricePoint
}

type priceClustererProvider func(points []providerPricePoint) priceClusterer

type priceFinder struct {
	log            *logan.Entry
	priceProviders []PriceProvider
	priceClustererProvider

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
		priceClustererProvider: func(points []providerPricePoint) priceClusterer {
			return newPriceClusterer(points)
		},
		maxPercentPriceDelta:        maxPercentPriceDelta,
		minNumberOfProvidersToAgree: minNumberOfProvidersToAgree,
	}

	return result, nil
}

func isPercentValid(percent int64) bool {
	return percent >= 0 && percent <= 100*amount.One
}

// TryFind - tries to find most recent Price Point which priceDelta is <= maxPercentPriceDelta
// for percent of providers >= minPercentOfNodeToParticipate
func (p *priceFinder) TryFind() (*provider.PricePoint, error) {
	allPoints, err := p.getAllProvidersPoints()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Points from all Providers")
	}

	sort.Sort(sortablePricePointsByTime(allPoints))
	clusterer := p.priceClustererProvider(allPoints)

	for i := range allPoints {
		// Cluster is a list of the closest PricePoints from Providers.
		cluster := clusterer.GetClusterForPoint(allPoints[i])
		candidate := median(cluster)

		if p.isMeetsReq(candidate, cluster) {
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

// RemoveOldPoints - removes points which were created before minAllowedTime
func (p *priceFinder) RemoveOldPoints(minAllowedTime time.Time) {
	for i := range p.priceProviders {
		p.priceProviders[i].RemoveOldPoints(minAllowedTime)
	}
}

// Median returns one of the provided points.
func median(points []providerPricePoint) providerPricePoint {
	sort.Sort(sortablePricePointsByPrice(points))
	// we can not select mean from two median points as point must exist
	return points[len(points)/2]
}

func (p *priceFinder) getAllProvidersPoints() ([]providerPricePoint, error) {
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

	return result, nil
}

// IsMeetsReq shows whether a candidate meets requirements.
// Returns true if cluster contains at least minNumberOfProvidersToAgree PricePoints, which agree with candidate.
func (p *priceFinder) isMeetsReq(candidate providerPricePoint, cluster []providerPricePoint) bool {
	if len(cluster) < p.minNumberOfProvidersToAgree {
		p.log.WithFields(logan.F{
			"min_number_of_providers": p.minNumberOfProvidersToAgree,
			"cluster_size":            len(cluster),
		}).Debug("Cluster too small - skipping.")
		return false
	}

	providersAgreed := 0
	for _, providerPricePoint := range cluster {
		if p.isMeetsPriceDeltaReq(candidate, providerPricePoint) {
			p.log.WithField("provider_price_point", providerPricePoint).Debug("Agreed.")
			providersAgreed++
		} else {
			p.log.WithField("provider_price_point", providerPricePoint).Debug("Disagreed.")
		}
	}

	return providersAgreed >= p.minNumberOfProvidersToAgree
}

func (p *priceFinder) isMeetsPriceDeltaReq(candidate, point providerPricePoint) bool {
	percentInDelta, isOverflow := candidate.GetPercentDeltaToMinPrice(point.PricePoint)
	if isOverflow {
		p.log.WithFields(logan.F{
			"candidate_price": candidate.Price,
			"point_price":     point.Price,
		}).Warn("Overflow on price delta calculation.")
		return false
	}

	return percentInDelta <= p.maxPercentPriceDelta
}
