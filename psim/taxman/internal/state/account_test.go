package state

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	exchange := AccountID("exchange_id")
	account := Account{
		Address: AccountID("account_id"),
		Balances: []*Balance{
			{
				Asset:      "XAAU",
				Address:    "xaau_balance",
				ExchangeID: exchange,
			},
		},
	}

	Convey("GetBalanceForAsset", t, func() {
		Convey("Balance with asset exists", func() {
			balance := account.GetBalanceForAsset(exchange, account.Balances[0].Asset)
			So(balance, ShouldNotBeNil)
		})
		Convey("Balance with asset does not exist", func() {
			balance := account.GetBalanceForAsset(exchange, "random_asset_code")
			So(balance, ShouldBeNil)
		})
	})
	Convey("MustGetBalanceForBalanceID", t, func() {
		Convey("Balance with balance ID exists", func() {
			balance := account.MustGetBalanceForBalanceID(account.Balances[0].Address)
			So(balance, ShouldNotBeNil)
		})
		Convey("Balance with asset does not exist", func() {
			So(func() {
				account.MustGetBalanceForBalanceID("random_balance_id")
			}, ShouldPanic)
		})
	})
	Convey("AddBalance", t, func() {
		Convey("Balance with such balance ID does not exist", func() {
			balance := Balance{
				Address: BalanceID("new_balance_id"),
			}
			err := account.AddBalance(balance)
			So(err, ShouldBeNil)
			actualBalance := account.MustGetBalanceForBalanceID(balance.Address)
			So(actualBalance, ShouldNotBeNil)
			assert.Equal(t, *actualBalance, balance)
		})
		Convey("Balance with such balance ID exists", func() {
			err := account.AddBalance(*account.Balances[0])
			So(err, ShouldNotBeNil)
		})
	})
}
