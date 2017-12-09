package txhandler

import (
	"math/rand"
	"testing"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	"gitlab.com/tokend/psim/psim/taxman/internal/txhandler/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/distributed_lab/logan/v3"
)

func TestLedgerChangesHandler(t *testing.T) {
	Convey("Given valid ledger changes handler", t, func() {
		statable := &mocks.Statable{}
		defer statable.AssertExpectations(t)
		handler := newLedgerChangesHandler(statable, logan.New())
		Convey("Failed to unmarshal", func() {
			err := handler.Handle(horizon.Transaction{})
			So(err, ShouldNotBeNil)
		})
		Convey("processBalanceCreated", func() {
			balance := validBalanceEntry()
			balanceChange := ledgerEntryChangeToXdrTxMeta([]xdr.LedgerEntryChange{
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
					Created: &xdr.LedgerEntry{
						Data: xdr.LedgerEntryData{
							Type:    xdr.LedgerEntryTypeBalance,
							Balance: balance,
						},
					},
				},
			})
			Convey("SpecialAccount", func() {
				statable.On("IsSpecialAccount", state.AccountID(balance.AccountId.Address())).Return(true).Once()
				Convey("Failed to add balance", func() {
					account := &state.Account{
						Balances: []*state.Balance{
							{
								Address: state.BalanceID(balance.BalanceId.AsString()),
							},
						},
					}
					statable.On("GetSpecialAccount", state.AccountID(balance.AccountId.Address())).Return(account).Once()
					err := handler.Handle(horizon.Transaction{
						ResultMetaXDR: balanceChange,
					})
					So(err, ShouldBeNil)
					So(len(account.Balances), ShouldEqual, 1)
				})
				Convey("Added balance", func() {
					account := &state.Account{}
					statable.On("GetSpecialAccount", state.AccountID(balance.AccountId.Address())).Return(account).Once()
					err := handler.Handle(horizon.Transaction{
						ResultMetaXDR: balanceChange,
					})
					So(err, ShouldBeNil)
					So(len(account.Balances), ShouldEqual, 1)
				})
			})
			Convey("Operational account", func() {
				statable.On("IsSpecialAccount", state.AccountID(balance.AccountId.Address())).Return(false).Once()
				operationalAccountID := state.AccountID(balance.AccountId.Address())
				statable.On("GetOperationalAccount").Return(operationalAccountID).Once()
				account := &state.Account{
					Balances: []*state.Balance{
						{
							Address: state.BalanceID(balance.BalanceId.AsString()),
						},
					},
				}
				statable.On("GetAccount", state.AccountID(balance.AccountId.Address())).Return(account).Once()
				err := handler.Handle(horizon.Transaction{
					ResultMetaXDR: balanceChange,
				})
				So(err, ShouldBeNil)
				So(len(account.Balances), ShouldEqual, 1)
			})
			Convey("Regular account", func() {
				statable.On("IsSpecialAccount", state.AccountID(balance.AccountId.Address())).Return(false).Once()
				Convey("Failed to add balance", func() {
					statable.On("GetOperationalAccount").Return(state.AccountID("operational_account")).Once()
					account := &state.Account{
						Balances: []*state.Balance{
							{
								Address: state.BalanceID(balance.BalanceId.AsString()),
							},
						},
					}
					statable.On("GetAccount", state.AccountID(balance.AccountId.Address())).Return(account).Once()
					err := handler.Handle(horizon.Transaction{
						ResultMetaXDR: balanceChange,
					})
					So(err, ShouldNotBeNil)
					So(len(account.Balances), ShouldEqual, 1)
				})
				Convey("Added balance", func() {
					account := &state.Account{}
					statable.On("GetAccount", state.AccountID(balance.AccountId.Address())).Return(account).Once()
					err := handler.Handle(horizon.Transaction{
						ResultMetaXDR: balanceChange,
					})
					So(err, ShouldBeNil)
					So(len(account.Balances), ShouldEqual, 1)
				})
			})
		})
		Convey("processAccountCreated", func() {
			account := validAccountEntry()
			accountChange := ledgerEntryChangeToXdrTxMeta([]xdr.LedgerEntryChange{
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
					Created: &xdr.LedgerEntry{
						Data: xdr.LedgerEntryData{
							Type:    xdr.LedgerEntryTypeAccount,
							Account: account,
						},
					},
				},
			})
			statable.On("AddAccount", state.Account{
				Address:          state.AccountID(account.AccountId.Address()),
				ShareForReferrer: int64(account.ShareForReferrer),
				Parent:           state.AccountID(account.Referrer.Address()),
			}).Once()
			err := handler.Handle(horizon.Transaction{
				ResultMetaXDR: accountChange,
			})
			So(err, ShouldBeNil)
		})
		Convey("processBalanceUpdated", func() {
			xdrBalance := validBalanceEntry()
			balanceChange := ledgerEntryChangeToXdrTxMeta([]xdr.LedgerEntryChange{
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
					Updated: &xdr.LedgerEntry{
						Data: xdr.LedgerEntryData{
							Type:    xdr.LedgerEntryTypeBalance,
							Balance: xdrBalance,
						},
					},
				},
			})
			Convey("SpecialAccount", func() {
				statable.On("IsSpecialAccount", state.AccountID(xdrBalance.AccountId.Address())).Return(true).Once()
				err := handler.Handle(horizon.Transaction{
					ResultMetaXDR: balanceChange,
				})
				So(err, ShouldBeNil)
			})
			Convey("RegularAccount", func() {
				statable.On("IsSpecialAccount", state.AccountID(xdrBalance.AccountId.Address())).Return(false).Once()
				balance := &state.Balance{

					Address: state.BalanceID(xdrBalance.BalanceId.AsString()),
				}
				account := &state.Account{
					Balances: []*state.Balance{balance},
				}
				statable.On("GetAccount", state.AccountID(xdrBalance.AccountId.Address())).Return(account).Once()
				err := handler.Handle(horizon.Transaction{
					ResultMetaXDR: balanceChange,
				})
				So(err, ShouldBeNil)
				So(len(account.Balances), ShouldEqual, 1)
				So(balance.FeesPaid, ShouldEqual, int64(xdrBalance.FeesPaid))
				So(balance.Amount, ShouldEqual, int64(xdrBalance.Amount))
			})
		})
		Convey("processAssetUpdated", func() {
			xdrAsset := xdr.AssetEntry{
				Code: xdr.AssetCode("XAAU"),
			}
			Convey("Add token", func() {
				tokenCode := xdr.AssetCode("XAAUT")
				xdrAsset.Token = &tokenCode
				assetChange := ledgerEntryChangeToXdrTxMeta([]xdr.LedgerEntryChange{
					{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
						Updated: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type:  xdr.LedgerEntryTypeAsset,
								Asset: &xdrAsset,
							},
						},
					},
				})
				statable.On("SetToken", state.AssetCode(xdrAsset.Code), state.AssetCode(*xdrAsset.Token)).Once()
				err := handler.Handle(horizon.Transaction{
					ResultMetaXDR: assetChange,
				})
				So(err, ShouldBeNil)
			})
			Convey("Update asset", func() {
				assetChange := ledgerEntryChangeToXdrTxMeta([]xdr.LedgerEntryChange{
					{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
						Updated: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type:  xdr.LedgerEntryTypeAsset,
								Asset: &xdrAsset,
							},
						},
					},
				})
				err := handler.Handle(horizon.Transaction{
					ResultMetaXDR: assetChange,
				})
				So(err, ShouldBeNil)
			})
		})
		Convey("LedgerEntry change not handled", func() {
			account := validAccountEntry()
			accountChange := ledgerEntryChangeToXdrTxMeta([]xdr.LedgerEntryChange{
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
					State: &xdr.LedgerEntry{
						Data: xdr.LedgerEntryData{
							Type:    xdr.LedgerEntryTypeAccount,
							Account: account,
						},
					},
				},
			})
			err := handler.Handle(horizon.Transaction{
				ResultMetaXDR: accountChange,
			})
			So(err, ShouldBeNil)

		})
	})
}

func validAccountEntry() *xdr.AccountEntry {
	account := xdr.AccountEntry{
		AccountType: xdr.AccountTypeGeneral,
	}
	So(account.AccountId.SetAddress("GBAOGSOXQHD25UVIVZ4XKQCYBHYYCFTWU2JO23PBBO5XA4JVZ3OMTOLN"), ShouldBeNil)
	account.Referrer = new(xdr.AccountId)
	So(account.Referrer.SetAddress("GBAOGSOXQHD25UVIVZ4XKQCYBHYYCFTWU2JO23PBBO5XA4JVZ3OMTOLN"), ShouldBeNil)
	account.ShareForReferrer = xdr.Int64(rand.Int63())
	return &account
}

func validBalanceEntry() *xdr.BalanceEntry {
	balanceID, err := horizon.ParseBalanceID("BBE674Z63XRPS2PSIQSYWBUBQFYMMMEHQJEUPDYJDPN3FPM6KXFQPF26")
	So(err, ShouldBeNil)
	balance := xdr.BalanceEntry{
		BalanceId: balanceID,
		Amount:    xdr.Int64(rand.Int63()),
		FeesPaid:  xdr.Int64(rand.Int63()),
	}
	So(balance.AccountId.SetAddress("GBAOGSOXQHD25UVIVZ4XKQCYBHYYCFTWU2JO23PBBO5XA4JVZ3OMTOLN"), ShouldBeNil)
	So(balance.Exchange.SetAddress("GBAOGSOXQHD25UVIVZ4XKQCYBHYYCFTWU2JO23PBBO5XA4JVZ3OMTOLN"), ShouldBeNil)
	return &balance
}

func ledgerEntryChangeToXdrTxMeta(changes []xdr.LedgerEntryChange) string {
	txMeta := xdr.TransactionMeta{
		Operations: &[]xdr.OperationMeta{
			{
				Changes: changes,
			},
		},
	}

	rawXDR, err := xdr.MarshalBase64(txMeta)
	So(err, ShouldBeNil)
	return string(rawXDR)
}
