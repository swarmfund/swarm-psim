package state

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	Convey("Given empty balance", t, func() {
		var emptyBalance Balance
		Convey("Set 0 fee (Demurrage) does not affect state", func() {
			emptyBalance.SetFeesPaid(0)
			assert.Equal(t, emptyBalance, Balance{})
		})
		Convey("Payout does not affect state", func() {
			emptyBalance.UpdateFeesAfterPayout()
			assert.Equal(t, emptyBalance, Balance{})
		})

	})
	Convey("Given balance with fees", t, func() {
		var balanceWithFees Balance
		feesPaid := rand.Int63()
		balanceWithFees.SetFeesPaid(feesPaid)
		Convey("Should share all the fees", func() {
			So(balanceWithFees.GetFeesToShare(), ShouldEqual, feesPaid)
		})
		Convey("After demurrage should share sum of fees set before demurrage and after", func() {
			checkDemurrage(&balanceWithFees, feesPaid)
			Convey("After payout fees to share should be 0", func() {
				checkPayout(&balanceWithFees)
			})
		})
		Convey("After payout fees to share should be 0", func() {
			checkPayout(&balanceWithFees)
			Convey("After demurrage should share sum of fees set before demurrage and after", func() {
				checkDemurrage(&balanceWithFees, feesPaid)
			})
		})
	})
}

func checkPayout(balance *Balance) {
	balance.UpdateFeesAfterPayout()
	So(balance.GetFeesToShare(), ShouldEqual, 0)
}

func checkDemurrage(balance *Balance, initialFeesToShare int64) {
	expectedFeesToShare := initialFeesToShare
	for i := 0; i < 2; i++ {
		// demurrage
		balance.SetFeesPaid(0)
		So(balance.GetFeesToShare(), ShouldEqual, expectedFeesToShare)
		// paid more fees after demurrage
		feesPaidAfterDemurrage := rand.Int63()
		balance.SetFeesPaid(feesPaidAfterDemurrage)
		expectedFeesToShare += feesPaidAfterDemurrage
		So(balance.GetFeesToShare(), ShouldEqual, expectedFeesToShare)
	}
}
