package snapshoter

import (
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"math/rand"
)

func TestPayoutBuilder(t *testing.T) {
	Convey("Given valid payout builder", t, func() {
		statable := &mockStatable{}
		defer statable.AssertExpectations(t)
		payoutBuilder := payoutBuilderImpl{
			state: statable,
		}

		masterAccountID := state.AccountID("master_account")
		statable.On("GetMasterAccount").Return(masterAccountID).Once()
		commissionAccount := state.Account{
			Address: state.AccountID("commission_id"),
			Balances: []*state.Balance{
				{
					Asset: state.AssetCode("XAAU"),
					ExchangeID: masterAccountID,
				},
			},
		}

		statable.On("GetCommissionAccount").Return(commissionAccount.Address).Once()
		specialAccounts := &state.Accounts{
			Accounts: map[state.AccountID]*state.Account{
				commissionAccount.Address: &commissionAccount,
			},
		}
		statable.On("GetSpecialAccounts").Return(specialAccounts).Once()
		ledger := rand.Int63()
		statable.On("GetLedger").Return(ledger).Once()

		Convey("payoutInfo has negative value", func() {
			err := payoutBuilder.BuildOperations(map[string]OperationSync{}, map[state.AccountID]map[state.AssetCode]int64{
				state.AccountID(""): {
					state.AssetCode(""): -123,
				},
			}, payoutTypeToken)
			So(err, ShouldNotBeNil)
		})
		Convey("nothing to share", func() {
			err := payoutBuilder.BuildOperations(map[string]OperationSync{}, map[state.AccountID]map[state.AssetCode]int64{
				state.AccountID(""): {
					state.AssetCode(""): 0,
				},
			}, payoutTypeToken)
			So(err, ShouldBeNil)
		})
		Convey("balance does not exist", func() {
			statable.On("GetAccount", state.AccountID("")).Return(&state.Account{}).Once()
			err := payoutBuilder.BuildOperations(map[string]OperationSync{}, map[state.AccountID]map[state.AssetCode]int64{
				state.AccountID(""): {
					state.AssetCode(""): rand.Int63() + 1,
				},
			}, payoutTypeToken)
			So(err, ShouldNotBeNil)
		})
		Convey("Given account with valid balance", func() {
			balanceID := state.BalanceID("account_XAAU_balance_ID")
			statable.On("GetAccount", state.AccountID("")).Return(&state.Account{
				Balances: []*state.Balance{
					{
						Address: balanceID,
						Asset: "XAAU",
						ExchangeID: masterAccountID,
					},
					{
						ExchangeID: masterAccountID,
					},
				},
			}).Once()
			Convey("Commission balance does not exist", func() {
				err := payoutBuilder.BuildOperations(map[string]OperationSync{}, map[state.AccountID]map[state.AssetCode]int64{
					state.AccountID(""): {
						state.AssetCode(""): rand.Int63() + 1,
					},
				}, payoutTypeToken)
				So(err, ShouldNotBeNil)
			})
			Convey("Payout with such reference already exist", func() {
				expectedReference := reference(payoutTypeToken, balanceID, ledger)
				err := payoutBuilder.BuildOperations(map[string]OperationSync{
					expectedReference: {},
				}, map[state.AccountID]map[state.AssetCode]int64{
					state.AccountID(""): {
						state.AssetCode("XAAU"): rand.Int63() + 1,
					},
				}, payoutTypeToken)
				So(err, ShouldNotBeNil)
			})
			Convey("Success", func() {
				result := map[string]OperationSync{}
				err := payoutBuilder.BuildOperations(result, map[state.AccountID]map[state.AssetCode]int64{
					state.AccountID(""): {
						state.AssetCode("XAAU"): rand.Int63() + 1,
					},
				}, payoutTypeToken)
				So(err, ShouldBeNil)
				So(len(result), ShouldEqual, 1)
			})
		})

	})
}
