package finder

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/prices/pricesetter/provider"
)

func TestPriceFinder(t *testing.T) {
	log := logan.New().WithField("service", "price_finder")
	Convey("Create price finder", t, func() {
		Convey("maxPercentPriceDelta < 0", func() {
			_, err := NewPriceFinder(log, nil, -1, 1)
			So(err, ShouldNotBeNil)
		})
		Convey("maxPercentPriceDelta > 100", func() {
			_, err := NewPriceFinder(log, nil, 101*amount.One, 1)
			So(err, ShouldNotBeNil)
		})
		Convey("minPercentOfProvidersToAgree < 0", func() {
			_, err := NewPriceFinder(log, nil, 1, -1)
			So(err, ShouldNotBeNil)
		})
		Convey("minPercentOfProvidersToAgree > 100", func() {
			_, err := NewPriceFinder(log, nil, 0, 101*amount.One)
			So(err, ShouldNotBeNil)
		})
		Convey("Empty priceProviders", func() {
			_, err := NewPriceFinder(log, nil, 0, 0)
			So(err, ShouldNotBeNil)
		})
	})
	Convey("Given valid PriceFinder", t, func() {
		mockedPriceProvider := MockPriceProvider{}
		defer mockedPriceProvider.AssertExpectations(t)
		maxPercentPriceDelta := int64(10 * amount.One)

		priceFinder, err := NewPriceFinder(log, []PriceProvider{&mockedPriceProvider}, maxPercentPriceDelta, 3)
		So(err, ShouldBeNil)

		Convey("Failed to get prices", func() {
			mockedPriceProvider.On("GetName").Return("mocked_provider").Once()
			expectedError := errors.New("Failed to get points")
			mockedPriceProvider.On("GetPoints").Return(nil, expectedError).Once()

			_, err := priceFinder.TryFind()
			So(errors.Cause(err), ShouldEqual, expectedError)
		})
		Convey("No PricePoints available", func() {
			mockedPriceProvider.On("GetPoints").Return(nil, nil).Once()
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

func testPriceFinding(t *testing.T, maxPercentPriceDelta int64, minPercentOfProvidersToAgree int,
	rates [][]provider.PricePoint, expectedResult *provider.PricePoint) {
	priceProviders := make([]PriceProvider, len(rates))
	for i := range priceProviders {
		priceProviderMock := MockPriceProvider{}
		//noinspection GoDeferInLoop
		defer priceProviderMock.AssertExpectations(t)
		priceProviderMock.On("GetPoints").Return(rates[i], nil).Once()
		priceProviderMock.On("GetName").Return(fmt.Sprintf("rates_provider_%d", i))
		priceProviders[i] = &priceProviderMock
	}

	priceFinder, err := NewPriceFinder(logan.New().WithField("service", "price_finder"), priceProviders,
		maxPercentPriceDelta, minPercentOfProvidersToAgree)
	So(err, ShouldBeNil)
	actualPoint, err := priceFinder.TryFind()
	So(err, ShouldBeNil)
	So(actualPoint, ShouldResemble, expectedResult)
}
