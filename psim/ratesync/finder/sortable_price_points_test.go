package finder

import (
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/swarmfund/psim/psim/ratesync/provider"
	"sort"
	"testing"
	"time"
)

func TestSortablePricePoints(t *testing.T) {

	Convey("Given array of price points", t, func() {
		input := []pricePoint{
			{PricePoint: provider.PricePoint{Time: time.Unix(3, 0), Price: 12}},
			{PricePoint: provider.PricePoint{Time: time.Unix(0, 0), Price: 15}},
			{PricePoint: provider.PricePoint{Time: time.Unix(1, 0), Price: 11}},
			{PricePoint: provider.PricePoint{Time: time.Unix(2, 0), Price: 13}},
		}

		Convey("Sort by time", func() {
			expected := []pricePoint{
				{PricePoint: provider.PricePoint{Time: time.Unix(3, 0), Price: 12}},
				{PricePoint: provider.PricePoint{Time: time.Unix(2, 0), Price: 13}},
				{PricePoint: provider.PricePoint{Time: time.Unix(1, 0), Price: 11}},
				{PricePoint: provider.PricePoint{Time: time.Unix(0, 0), Price: 15}},
			}

			sort.Sort(sortablePricePointsByTime(input))
			for i := range input {
				So(input[i], ShouldResemble, expected[i])
			}
		})
		Convey("Sort by price", func() {
			expected := []pricePoint{
				{PricePoint: provider.PricePoint{Time: time.Unix(0, 0), Price: 15}},
				{PricePoint: provider.PricePoint{Time: time.Unix(2, 0), Price: 13}},
				{PricePoint: provider.PricePoint{Time: time.Unix(3, 0), Price: 12}},
				{PricePoint: provider.PricePoint{Time: time.Unix(1, 0), Price: 11}},
			}

			sort.Sort(sortablePricePointsByPrice(input))
			for i := range input {
				So(input[i], ShouldResemble, expected[i])
			}
		})

	})
}
