package finder

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"time"
	"fmt"
)

func TestPriceFinder(t *testing.T) {
	log := logan.New().WithField("service", "price_finder")
	Convey("Create price finder", t, func() {
		Convey("maxPercentPriceDelta < 0", func() {
			_, err := newPriceFinder(log, nil, -1, 1, nil)
			So(err, ShouldNotBeNil)
		})
		Convey("maxPercentPriceDelta > 100", func() {
			_, err := newPriceFinder(log, nil, 101*amount.One, 1, nil)
			So(err, ShouldNotBeNil)
		})
		Convey("minPercentOfProvidersToAgree < 0", func() {
			_, err := newPriceFinder(log, nil, 1, -1, nil)
			So(err, ShouldNotBeNil)
		})
		Convey("minPercentOfProvidersToAgree > 100", func() {
			_, err := newPriceFinder(log, nil, 0, 101*amount.One, nil)
			So(err, ShouldNotBeNil)
		})
		Convey("Empty ratesProviders", func() {
			_, err := newPriceFinder(log, nil, 0, 0, nil)
			So(err, ShouldNotBeNil)
		})
	})
	Convey("Given valid price finder", t, func() {
		mockedRatesProvider := mockRatesProvider{}
		defer mockedRatesProvider.AssertExpectations(t)
		maxPercentPriceDelta := int64(10 * amount.One)
		minPercentOfProviderToAgree := int64(51 * amount.One)
		priceClustererMock := mockPriceClusterer{}
		priceFinder, err := newPriceFinder(log, []ratesProvider{&mockedRatesProvider}, maxPercentPriceDelta,
			minPercentOfProviderToAgree, func(points []pricePoint) priceClusterer {
				return &priceClustererMock
			})
		So(err, ShouldBeNil)
		Convey("Failed to get prices", func() {
			mockedRatesProvider.On("GetName").Return("mocked_provider").Once()
			expectedError := errors.New("Failed to get points")
			mockedRatesProvider.On("GetPoints").Return(nil, expectedError).Once()
			_, err := priceFinder.TryFind()
			So(errors.Cause(err), ShouldEqual, expectedError)
		})
		Convey("No price points available", func() {
			mockedRatesProvider.On("GetPoints").Return(nil, nil).Once()
			result, err := priceFinder.TryFind()
			So(err, ShouldBeNil)
			So(result, ShouldBeNil)
		})
	})
	Convey("Price finding", t, func() {
		Convey("One provider, just select latest point", func() {
			expectedResult := &provider.PricePoint{Price: 20, Time: time.Unix(11, 0)}
			testPriceFinding(t, 10*amount.One, 51*amount.One, [][]provider.PricePoint{
				{
					{20, time.Unix(11, 0)}, {101, time.Unix(9, 0)}, {100, time.Unix(8, 0)},
				},
			}, expectedResult)
		})
		Convey("Three providers acting normally", func() {
			expectedResult := &provider.PricePoint{Price: 20, Time: time.Unix(11, 0)}
			testPriceFinding(t, 10*amount.One, 51*amount.One, [][]provider.PricePoint{
				{
					{20, time.Unix(11, 0)}, {101, time.Unix(9, 0)}, {100, time.Unix(8, 0)},
					{20, time.Unix(11, 0)}, {101, time.Unix(9, 0)}, {100, time.Unix(8, 0)},
					{21, time.Unix(11, 0)}, {101, time.Unix(9, 0)}, {100, time.Unix(8, 0)},
				},
			}, expectedResult)
		})
		Convey("Three providers - one provider delayed info", func() {
			expectedResult := &provider.PricePoint{Price: 101, Time: time.Unix(11, 0)}
			testPriceFinding(t, 10*amount.One, 51*amount.One, [][]provider.PricePoint{
				{
					{101, time.Unix(11, 0)}, {80, time.Unix(9, 0)}, {60, time.Unix(8, 0)},
					{20, time.Unix(11, 0)}, {101, time.Unix(9, 0)}, {81, time.Unix(8, 0)},
					{21, time.Unix(11, 0)}, {101, time.Unix(9, 0)}, {82, time.Unix(8, 0)},
				},
			}, expectedResult)
		})
	})
}

func testPriceFinding(t *testing.T, maxPercentPriceDelta int64, minPercentOfProvidersToAgree int64,
	rates [][]provider.PricePoint, expectedResult *provider.PricePoint) {
	ratesProviders := make([]ratesProvider, len(rates))
	for i := range ratesProviders {
		ratesProviderMock := mockRatesProvider{}
		//noinspection GoDeferInLoop
		defer ratesProviderMock.AssertExpectations(t)
		ratesProviderMock.On("GetPoints").Return(rates[i], nil).Once()
		ratesProviderMock.On("GetName").Return(fmt.Sprintf("rates_provider_%d", i))
		ratesProviders[i] = &ratesProviderMock
	}

	priceFinder, err := newPriceFinder(logan.New().WithField("service", "price_finder"), ratesProviders,
		maxPercentPriceDelta, minPercentOfProvidersToAgree, func(points []pricePoint) priceClusterer {
			return newPriceClusterer(points)
		})
	So(err, ShouldBeNil)
	actualPoint, err := priceFinder.TryFind()
	So(err, ShouldBeNil)
	So(actualPoint, ShouldResemble, expectedResult)
}
