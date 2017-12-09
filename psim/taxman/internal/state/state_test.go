package state

import (
	"math"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	Convey("Given valid state", t, func() {
		state := NewState()
		Convey("Given valid ledger", func() {
			expectedLedger := rand.Int63()
			state.SetLedger(expectedLedger)
			Convey("Can get", func() {
				actualLedger := state.GetLedger()
				So(actualLedger, ShouldEqual, expectedLedger)
			})
		})
		Convey("Given valid token", func() {
			expectedAsset := AssetCode("XAAU")
			expectedAssetToken := AssetCode("XAAUT")
			state.SetToken(expectedAsset, expectedAssetToken)
			Convey("Can get asset by token", func() {
				actualAsset := state.GetAssetByToken(expectedAssetToken)
				So(actualAsset, ShouldEqual, expectedAsset)
			})
			Convey("panics if asset for token does not exist", func() {
				So(func() {
					state.GetAssetByToken("token_and_asset_does_not_exists")
				}, ShouldPanic)
			})
			Convey("isToken", func() {
				So(state.isToken(expectedAssetToken), ShouldBeTrue)
				So(state.isToken(expectedAsset), ShouldBeFalse)
			})
		})
		Convey("GetTotalFeesToShare", func() {
			state.SetToken("XAAU", "XAAUT")
			state.SetToken("XAAG", "XAAGT")
			account1 := Account{
				Address: AccountID("account_1"),
				Balances: []*Balance{
					{
						Asset:    AssetCode("XAAU"),
						FeesPaid: 100,
					},
					{
						Asset:    AssetCode("XAAUT"),
						FeesPaid: 120,
						Amount:   1,
					},
					{
						Asset:    AssetCode("XAAGT"),
						FeesPaid: 500,
						Amount:   120,
					},
				},
			}
			state.AddAccount(account1)
			account2 := Account{
				Address: AccountID("account_2"),
				Balances: []*Balance{
					{
						Asset:    AssetCode("XAAG"),
						FeesPaid: 300,
						Amount:   500,
					},
					{
						Asset:    AssetCode("XAAUT"),
						FeesPaid: 320,
						Amount:   505,
					},
				},
			}
			state.AddAccount(account2)

			actualFeesToShare := state.GetTotalFeesToShare()
			So(len(actualFeesToShare), ShouldEqual, 4)
			So(actualFeesToShare["XAAU"], ShouldEqual, 100)
			So(actualFeesToShare["XAAG"], ShouldEqual, 300)
			So(actualFeesToShare["XAAUT"], ShouldEqual, 440)
			So(actualFeesToShare["XAAGT"], ShouldEqual, 500)
			Convey("GetTotalFeesToShare - overflow", func() {
				account2.GetBalanceForAsset(AccountID(""), AssetCode("XAAUT")).FeesPaid = math.MaxInt64
				So(func() {
					state.GetTotalFeesToShare()
				}, ShouldPanic)
			})
			Convey("PayoutCompleted", func() {
				state.PayoutCompleted()
				actualFeesToShare = state.GetTotalFeesToShare()
				for _, actualFeeToShare := range actualFeesToShare {
					So(actualFeeToShare, ShouldEqual, 0)
				}
			})
			Convey("TokenBalances", func() {
				actualTokensBalances := state.TokenBalances()
				actualTokensBalancesSlice := []Balance{}
				for actualTokenBalance := range actualTokensBalances {
					actualTokensBalancesSlice = append(actualTokensBalancesSlice, *actualTokenBalance)
				}
				So(len(actualTokensBalancesSlice), ShouldEqual, 3)
				So(actualTokensBalancesSlice, ShouldContain, *state.GetMainBalanceForAsset(account1.Address, "XAAUT"))
				So(actualTokensBalancesSlice, ShouldContain, *state.GetMainBalanceForAsset(account1.Address, "XAAGT"))
				So(actualTokensBalancesSlice, ShouldContain, *state.GetMainBalanceForAsset(account2.Address, "XAAUT"))
			})
			Convey("GetTotalTokensAmount", func() {
				actualTotalTokensAmount := state.GetTotalTokensAmount()
				So(len(actualTotalTokensAmount), ShouldEqual, 2)
				So(actualTotalTokensAmount["XAAUT"], ShouldEqual, 506)
				So(actualTotalTokensAmount["XAAGT"], ShouldEqual, 120)
				Convey("Panic on overflow", func() {
					actualAccount2 := state.GetAccount(account2.Address)
					actualBalance2 := actualAccount2.GetBalanceForAsset(AccountID(""), "XAAUT")
					actualBalance2.Amount = math.MaxInt64
					So(func() {
						state.GetTotalTokensAmount()
					}, ShouldPanic)
				})
			})
		})
		Convey("Special accounts", func() {
			state.SetSpecialAccounts(AccountID("master"), AccountID("storage"), AccountID("commission"))
			So(state.GetMasterAccount(), ShouldEqual, AccountID("master"))
			So(state.GetStorageFeeAccount(), ShouldEqual, AccountID("storage"))
			So(state.GetCommissionAccount(), ShouldEqual, AccountID("commission"))
			specialAccount := Account{
				Address: AccountID("special_account"),
			}
			state.GetSpecialAccounts().AddAccount(specialAccount)
			assert.Equal(t, *state.GetSpecialAccounts().GetAccount(specialAccount.Address), specialAccount)
			assert.Equal(t, *state.GetSpecialAccount(specialAccount.Address), specialAccount)
			So(state.IsSpecialAccount(specialAccount.Address), ShouldBeTrue)
			notSpecialAccount := Account{
				Address: "not_special_account",
			}
			state.AddAccount(notSpecialAccount)
			So(state.IsSpecialAccount(notSpecialAccount.Address), ShouldBeFalse)
		})
		Convey("Operational account", func() {
			state.SetOperationalAccount("operational_account")
			So(state.GetOperationalAccount(), ShouldEqual, "operational_account")
		})
		Convey("Payout period", func() {
			So(state.GetPayoutPeriod(), ShouldBeNil)
			expectedPayout := time.Duration(rand.Int63()) * time.Second
			state.SetPayoutPeriod(&expectedPayout)
			So(*state.GetPayoutPeriod(), ShouldEqual, expectedPayout)
		})
	})
}
