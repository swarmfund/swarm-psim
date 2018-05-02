package finder

import (
	"sort"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/prices/types"
	"gitlab.com/tokend/go/amount"
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
	maxPercentPriceDelta int64
	minProvidersToAgree  int
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
		return nil, errors.New("minProvidersToAgree must be > 0")
	}

	if len(priceProviders) == 0 {
		return nil, errors.New("Unexpected number of PriceProviders.")
	}

	result := &priceFinder{
		log:            log.WithField("service", "price_finder"),
		priceProviders: priceProviders,
		priceClusterizerProvider: func(points []providerPricePoint) priceClusterizer {
			return newPriceClusterizer(points)
		},
		maxPercentPriceDelta: maxPercentPriceDelta,
		minProvidersToAgree:  minNumberOfProvidersToAgree,
	}

	return result, nil
}

func isPercentValid(percent int64) bool {
	return percent >= 0 && percent <= 100*amount.One
}

// TryFind - tries to find most recent PricePoint which priceDelta is <= maxPercentPriceDelta
// for percent of providers >= minPercentOfNodeToParticipate
func (f *priceFinder) TryFind() (*types.PricePoint, error) {
	f.log.Debug("Tyring to find a PricePoint.")

	allPoints := f.getAllProvidersPoints()

	sort.Sort(sortablePricePointsByTime(allPoints))
	clusterizer := f.priceClusterizerProvider(allPoints)

	for i := range allPoints {
		cluster := clusterizer.GetClusterForPoint(allPoints[i].PricePoint)
		candidate := median(cluster)
		fields := logan.F{
			"candidate": candidate,
			"cluster":   cluster,
		}

		verifyErr := f.verifyCandidate(candidate.Price, cluster)
		if verifyErr != nil {
			f.log.WithFields(fields).WithError(verifyErr).Debug("Candidate is not approved over Cluster.")
			continue
		}

		return &candidate.PricePoint, nil
	}

	return nil, nil
}

// VerifyPrice returns true if median of the Cluster built over all existing Points
// meets delta requirements with the provided point.
func (f *priceFinder) VerifyPrice(price int64) error {
	allPoints := f.getAllProvidersPoints()
	clusterizer := f.priceClusterizerProvider(allPoints)

	cluster := clusterizer.GetClusterForPoint(types.PricePoint{})

	priceVerifyErr := f.verifyCandidate(price, cluster)
	if priceVerifyErr != nil {
		return priceVerifyErr
	}

	m := median(cluster)
	f.RemoveOldPoints(m.Time)
	return nil
}

// RemoveOldPoints - removes points which were created before minAllowedTime
func (f *priceFinder) RemoveOldPoints(minAllowedTime time.Time) {
	for i := range f.priceProviders {
		f.priceProviders[i].RemoveOldPoints(minAllowedTime)
	}
}

func (f *priceFinder) getAllProvidersPoints() []providerPricePoint {
	var result []providerPricePoint

	for _, priceProvider := range f.priceProviders {
		pricePoints := priceProvider.GetPoints()

		for i := range pricePoints {
			result = append(result, providerPricePoint{
				ProviderID: priceProvider.GetName(),
				PricePoint: pricePoints[i],
			})
		}
	}

	return result
}

// Median returns one of the provided points - median by Price.
func median(points []providerPricePoint) providerPricePoint {
	sort.Sort(sortablePricePointsByPrice(points))
	// we can not select mean from two median points as point must exist
	return points[len(points)/2]
}

// VerifyCandidate shows whether provided candidatePrice meets requirements over provided cluster.
// Returns nil if cluster contains at least minProvidersToAgree PricePoints, which agree with candidate.
func (f *priceFinder) verifyCandidate(candidatePrice int64, cluster []providerPricePoint) (verifyErr error) {
	if len(cluster) < f.minProvidersToAgree {
		clusterProviders := make([]string, 0)
		for _, provider := range cluster {
			clusterProviders = append(clusterProviders, provider.ProviderID)
		}

		return errors.From(errors.New("Cluster too small, don't even trying to evaluate distances."), logan.F{
			"min_providers_to_agree": f.minProvidersToAgree,
			"cluster_size":           len(cluster),
			"cluster_providers":      clusterProviders,
		})
	}

	providersAgreed := 0
	disagreedProviders := make([]string, 0)
	for _, providerPricePoint := range cluster {
		if f.meetsPriceDeltaReq(candidatePrice, providerPricePoint.PricePoint.Price) {
			// Provider agreed with the candidate - diff between the providerPricePoint and the candidate isn't too big.
			f.log.WithField("provider_price_point", providerPricePoint).Debug("Provider agreed.")
			providersAgreed++
		} else {
			disagreedProviders = append(disagreedProviders, providerPricePoint.ProviderID)
			f.log.WithField("provider_price_point", providerPricePoint).Debug("Provider disagreed.")
		}
	}

	if providersAgreed < f.minProvidersToAgree {
		return errors.From(errors.New("Too few Providers agreed."), logan.F{
			"providers_agreed":       providersAgreed,
			"min_providers_to_agree": f.minProvidersToAgree,
			"disagreed_providers":    disagreedProviders,
		})
	}

	return nil
}

func (f *priceFinder) meetsPriceDeltaReq(price1, price2 int64) bool {
	percentInDelta, isOverflow := types.GetPercentDeltaToMinPrice(price1, price2)
	if isOverflow {
		f.log.WithFields(logan.F{
			"price_1": price1,
			"price_2": price2,
		}).Warn("Overflow on price delta calculation.")
		return false
	}

	return percentInDelta <= f.maxPercentPriceDelta
}
