package snapshoter

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

func TestReferralsShareProvider(t *testing.T) {
	Convey("Given valid referralShareProviderImpl", t, func() {
		statable := &mockStatable{}
		defer statable.AssertExpectations(t)
		referralShareProvider := referralShareProviderImpl{
			state: statable,
		}
		Convey("Invalid share for referrer - negative", func() {
			children := getChildren(&state.Account{
				ShareForReferrer: -100,
			})
			statable.On("GetChildren").Return(children).Once()
			_, err := referralShareProvider.GetReferralSharePayout()
			So(err, ShouldNotBeNil)
		})
		Convey("Invalid share for referrer - > 100%", func() {
			children := getChildren(&state.Account{
				ShareForReferrer: 101 * amount.One,
			})
			statable.On("GetChildren").Return(children).Once()
			_, err := referralShareProvider.GetReferralSharePayout()
			So(err, ShouldNotBeNil)
		})
		Convey("Share to parent overflow", func() {
			children := getChildren(&state.Account{
				ShareForReferrer: 100 * amount.One,
				Balances: []*state.Balance{
					{
						FeesPaid: math.MaxInt64,
					},
					{
						FeesPaid: math.MaxInt64,
					},
				},
			})
			statable.On("GetChildren").Return(children).Once()
			_, err := referralShareProvider.GetReferralSharePayout()
			So(err, ShouldNotBeNil)
		})
		Convey("Success", func() {
			children := getChildren(&state.Account{
				ShareForReferrer: 10 * amount.One,
				Balances: []*state.Balance{
					{
						FeesPaid: 256,
					},
					{
						FeesPaid: 512,
					},
				},
			})
			statable.On("GetChildren").Return(children).Once()
			referralsShare, err := referralShareProvider.GetReferralSharePayout()
			So(err, ShouldBeNil)
			So(len(referralsShare), ShouldEqual, 1)
			So(len(referralsShare[state.AccountID("")]), ShouldEqual, 1)
			So(referralsShare[state.AccountID("")][state.AssetCode("")], ShouldEqual, 76)
		})
	})
}

func getChildren(accounts ...*state.Account) chan *state.Account {
	result := make(chan *state.Account)
	go func() {
		for i := range accounts {
			result <- accounts[i]
		}

		close(result)
	}()

	return result
}
