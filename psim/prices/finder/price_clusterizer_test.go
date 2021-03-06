package finder

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
)

func TestPriceClusterizer(t *testing.T) {
	Convey("Panic when trying to calc distance to point from same provider", t, func() {
		So(func() {
			calcDistanceAbs(providerPricePoint{}, providerPricePoint{})
		}, ShouldPanic)
	})
	Convey("GetClusterForPoint", t, func() {
		Convey("Providers are always unique", func() {
			totalNumberOfPoints := rand.Int31n(1000)
			providers := []string{"p1", "p2", "p3", "p4", "p5"}
			input := make([]providerPricePoint, totalNumberOfPoints)
			for i := range input {
				input[i] = providerPricePoint{
					ProviderID: providers[rand.Intn(len(providers))],
					PricePoint: types.PricePoint{
						Price: rand.Int63(),
					},
				}
			}

			clusterizer := newPriceClusterizer(input)
			result := clusterizer.GetClusterForPoint(input[rand.Intn(len(input))])
			So(len(result), ShouldEqual, len(providers))

			usedProviders := map[string]bool{}
			for i := range result {
				_, exists := usedProviders[result[i].ProviderID]
				So(exists, ShouldBeFalse)
				usedProviders[result[i].ProviderID] = true
			}
		})
	})
}
