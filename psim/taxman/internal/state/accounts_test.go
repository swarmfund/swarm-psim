package state

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestAccounts(t *testing.T) {
	Convey("Given accounts with account", t, func() {
		accounts := NewAccounts()
		account := Account{
			Address: AccountID("account_id"),
			Balances: []*Balance{
				{
					Account: AccountID("account_id"),
					Asset:   "asset",
					Address: BalanceID("balance_id"),
					Amount:  rand.Int63(),
				},
			},
			Parent:           AccountID("child_account_id"),
			ShareForReferrer: rand.Int63(),
		}
		accounts.AddAccount(account)

		Convey("Can't add same account twice", func() {
			So(func() {
				accounts.AddAccount(account)
			}, ShouldPanic)
		})
		Convey("Can get existing account", func() {
			actualAccount := accounts.GetAccount(account.Address)
			So(actualAccount, ShouldNotBeNil)
			assert.Equal(t, *actualAccount, account)
		})
		Convey("Can't get non existing account", func() {
			So(func() {
				accounts.GetAccount(AccountID("account_does_not_exist"))
			}, ShouldPanic)
		})
		Convey("Exists", func() {
			So(accounts.Exists(account.Address), ShouldBeTrue)
			So(accounts.Exists("account_does_not_exist"), ShouldBeFalse)
		})
	})
	Convey("GetChildren", t, func() {
		accounts := NewAccounts()
		childAccount1 := Account{
			Address: AccountID("account_id"),
			Balances: []*Balance{
				{
					Account: AccountID("account_id"),
					Asset:   "asset",
					Address: BalanceID("balance_id"),
					Amount:  rand.Int63(),
				},
			},
			Parent:           AccountID("child_account_id"),
			ShareForReferrer: rand.Int63(),
		}
		accounts.AddAccount(childAccount1)
		childAccount2 := childAccount1
		childAccount2.Address = AccountID("account_id_2")
		accounts.AddAccount(childAccount2)

		orphanAccount := childAccount2
		orphanAccount.Address = AccountID("orphan")
		orphanAccount.ShareForReferrer = 0
		orphanAccount.Parent = AccountID("")
		accounts.AddAccount(orphanAccount)

		actualChildren := accounts.GetChildren()
		actualChildrenSlice := []Account{}
		for actualChild := range actualChildren {
			actualChildrenSlice = append(actualChildrenSlice, *actualChild)
		}
		So(actualChildrenSlice, ShouldContain, childAccount1)
		So(actualChildrenSlice, ShouldContain, childAccount2)
		So(actualChildrenSlice, ShouldNotContain, orphanAccount)
	})
}
