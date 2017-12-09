package snapshoter

import (
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTokensShareProvider(t *testing.T) {
	Convey("Given valid tokenShareProviderImpl", t, func() {
		statable := &mockStatable{}
		defer statable.AssertExpectations(t)
		tokenShareProvider := tokenShareProviderImpl{
			state: statable,
		}
		Convey("Invalid total tokens amount", func() {
			statable.On("GetTotalTokensAmount").Return(map[state.AssetCode]int64{
				"XAAUT": 0,
				"USDT":  10000,
				"ASDT":  -1000,
				"XAAGT": 123,
			}).Once()
			_, err := tokenShareProvider.GetFeesToShareToTokenHolders()
			So(err, ShouldNotBeNil)
		})
		Convey("Given valid total tokens amount", func() {
			statable.On("GetTotalTokensAmount").Return(map[state.AssetCode]int64{
				"XAAUT": 0,
				"USDT":  10000,
				"ASDT":  0,
			}).Once()
			Convey("Invalid total fees to share", func() {
				statable.On("GetTotalFeesToShare").Return(map[state.AssetCode]int64{
					"XAAUT": 0,
					"USDT":  10000,
					"ASDT":  -123,
					"XAAGT": 123,
				}).Once()
				_, err := tokenShareProvider.GetFeesToShareToTokenHolders()
				So(err, ShouldNotBeNil)
			})
			Convey("Given valid total fees to share", func() {
				statable.On("GetTotalFeesToShare").Return(map[state.AssetCode]int64{
					"XAAUT": 0,
					"USDT":  0,
					"USD":   50001,
					"XAAGT": 123,
				}).Once()
				Convey("token balance exceeds total amount", func() {
					balances := make(chan *state.Balance, 1)
					balances <- &state.Balance{
						Asset:  state.AssetCode("USDT"),
						Amount: 10001,
					}
					statable.On("TokenBalances").Return(balances).Once()
					_, err := tokenShareProvider.GetFeesToShareToTokenHolders()
					So(err, ShouldNotBeNil)
				})
				Convey("Success", func() {
					balances := make(chan *state.Balance, 3)
					balances <- &state.Balance{
						Asset:  state.AssetCode("USDT"),
						Amount: 2000,
					}
					balances <- &state.Balance{
						Asset:  state.AssetCode("USDT"),
						Amount: 3000,
					}
					balances <- &state.Balance{
						Asset:  state.AssetCode("XAAUT"),
						Amount: 3000,
					}
					close(balances)
					statable.On("TokenBalances").Return(balances).Once()
					statable.On("GetAssetByToken", state.AssetCode("USDT")).Return(state.AssetCode("USD")).Twice()
					result, err := tokenShareProvider.GetFeesToShareToTokenHolders()
					So(err, ShouldBeNil)
					So(len(result), ShouldEqual, 1)
					So(result[state.AccountID("")][state.AssetCode("USD")], ShouldEqual, 25000)
				})
			})
		})
	})
}
