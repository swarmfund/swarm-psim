package finder

import (
	"sort"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
)

func TestSortablePricePoints(t *testing.T) {

	Convey("Given array of price points", t, func() {
		input := []providerPricePoint{
			{PricePoint: types.PricePoint{Time: time.Unix(3, 0), Price: 12}},
			{PricePoint: types.PricePoint{Time: time.Unix(0, 0), Price: 15}},
			{PricePoint: types.PricePoint{Time: time.Unix(1, 0), Price: 11}},
			{PricePoint: types.PricePoint{Time: time.Unix(2, 0), Price: 13}},
		}

		Convey("Sort by time", func() {
			expected := []providerPricePoint{
				{PricePoint: types.PricePoint{Time: time.Unix(3, 0), Price: 12}},
				{PricePoint: types.PricePoint{Time: time.Unix(2, 0), Price: 13}},
				{PricePoint: types.PricePoint{Time: time.Unix(1, 0), Price: 11}},
				{PricePoint: types.PricePoint{Time: time.Unix(0, 0), Price: 15}},
			}

			sort.Sort(sortablePricePointsByTime(input))
			for i := range input {
				So(input[i], ShouldResemble, expected[i])
			}
		})
		Convey("Sort by price", func() {
			expected := []providerPricePoint{
				{PricePoint: types.PricePoint{Time: time.Unix(0, 0), Price: 15}},
				{PricePoint: types.PricePoint{Time: time.Unix(2, 0), Price: 13}},
				{PricePoint: types.PricePoint{Time: time.Unix(3, 0), Price: 12}},
				{PricePoint: types.PricePoint{Time: time.Unix(1, 0), Price: 11}},
			}

			sort.Sort(sortablePricePointsByPrice(input))
			for i := range input {
				So(input[i], ShouldResemble, expected[i])
			}
		})

	})
}
