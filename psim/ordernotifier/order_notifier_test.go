package ordernotifier

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/ordernotifier/mocks"
	"testing"
	"gitlab.com/swarmfund/go/xdr"
	"fmt"
)

func TestOrderNotifierHelper(t *testing.T) {
	ctx := context.Background()
	transactionConnector := &mocks.TransactionConnector{}
	emailSender := &mocks.NotificatorConnector{}
	logger := logan.New().WithField("service", "order_notifier")
	userConnector := &mocks.UserConnector{}
	saleConnector := &mocks.SaleConnector{}
	config := Config{
		EmailsConfig: EmailsConfig{
			Subject:            "IO cancelled",
			RequestType:        10,
			RequestTokenSuffix: "_io_cancelled",
			TemplateName:       "io_cancelled.html",
			TemplateLinkURL:    "https://invest.swarm.fund",
		},
	}

	s := New(config, transactionConnector, emailSender, logger, userConnector, saleConnector, make(chan horizon.CheckSaleStateResponse, 0))

	Convey("Ineligible effect", t, func() {
		checkSaleStateOp := horizon.CheckSaleState{
			SaleID:        1,
			Effect:        xdr.CheckSaleStateEffectClosed.String(),
			TransactionID: "some_transaction_id",
		}
		err := s.processCheckSaleStateOperation(ctx, checkSaleStateOp)
		So(err, ShouldBeNil)
	})
	Convey("Transaction doesn't exist", t, func() {
		defer transactionConnector.AssertExpectations(t)

		checkSaleStateOp := horizon.CheckSaleState{
			SaleID:        1,
			Effect:        xdr.CheckSaleStateEffectUpdated.String(),
			TransactionID: "some_transaction_id",
		}

		transactionConnector.On("TransactionByID", checkSaleStateOp.TransactionID).Return(nil, nil).Once()

		err := s.processCheckSaleStateOperation(ctx, checkSaleStateOp)
		So(err, ShouldBeNil)
	})
	Convey("Wrong change type", t, func() {
		change := xdr.LedgerEntryChange{
			Type: xdr.LedgerEntryChangeTypeCreated,
		}

		_, err := s.processLedgerEntry(ctx, change, 1)
		So(err, ShouldBeNil)
	})
	Convey("Wrong entry ", t, func() {
		change := xdr.LedgerEntryChange{
			Type: xdr.LedgerEntryChangeTypeRemoved,
			Removed: &xdr.LedgerKey{
				Type: xdr.LedgerEntryTypeAccount,
			},
		}

		_, err := s.processLedgerEntry(ctx, change, 1)
		So(err, ShouldBeNil)
	})
	Convey("User doesn't exist", t, func() {
		var saleID uint64 = 1
		var accountID xdr.AccountId
		err := accountID.SetAddress("GDDDAXSERVQRFM3STE65GNLWJ7A6QHOHDGJLSGFZBKETSPGR3RHI5TXT")
		offer := xdr.LedgerKeyOffer{
			OwnerId: accountID,
			OfferId: 1,
		}

		userConnector.On("User", offer.OwnerId.Address()).Return(nil, nil).Once()

		emailUnit, err := s.processCancelledOrder(ctx, offer, saleID)
		userConnector.AssertExpectations(t)
		So(emailUnit, ShouldBeNil)
		So(err, ShouldBeNil)

		Convey("User exists", func() {
			user := &horizon.User{
				Type: "General",
				ID:   "User id",
				Attributes: horizon.UserAttributes{
					Email: "test_mail@gmail.com",
				},
			}

			userConnector.On("User", offer.OwnerId.Address()).Return(user, nil).Once()
			defer userConnector.AssertExpectations(t)

			Convey("Valid unique token", func() {
				expectedUniqueToken := fmt.Sprintf("%s:%d:%d:%s", user.Attributes.Email, offer.OfferId, saleID, s.config.RequestTokenSuffix)
				actualUniqueToken := s.buildUniqueToken(user.Attributes.Email, uint64(offer.OfferId), uint64(saleID))
				So(expectedUniqueToken, ShouldEqual, actualUniqueToken)

				Convey("Sale doesn't exist", func() {
					saleConnector.On("SaleByID", saleID).Return(nil, nil).Once()
					defer saleConnector.AssertExpectations(t)

					emailUnit, err = s.processCancelledOrder(ctx, offer, saleID)
					So(emailUnit, ShouldBeNil)
					So(err, ShouldBeNil)

					Convey("Sale exists", func() {
						sale := &horizon.Sale{
							ID: saleID,
							Details: horizon.SaleDetails{
								Name: "Test sale",
							},
						}

						msg, err := s.buildEmailMessage(sale.Name())
						So(err, ShouldBeNil)

						userConnector.On("User", offer.OwnerId.Address()).Return(user, nil).Once()
						saleConnector.On("SaleByID", saleID).Return(sale, nil).Once()

						emailUnit, err = s.processCancelledOrder(ctx, offer, saleID)
						So(err, ShouldBeNil)
						So(emailUnit.UniqueToken, ShouldEqual, expectedUniqueToken)
						So(emailUnit.Payload.Destination, ShouldEqual, user.Attributes.Email)
						So(emailUnit.Payload.Subject, ShouldEqual, s.config.Subject)
						So(emailUnit.Payload.Message, ShouldEqual, msg)
					})
				})
			})
		})
	})
}
