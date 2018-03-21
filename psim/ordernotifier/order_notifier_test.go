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

	var saleID uint64 = 1
	var accountID xdr.AccountId
	err := accountID.SetAddress("GDDDAXSERVQRFM3STE65GNLWJ7A6QHOHDGJLSGFZBKETSPGR3RHI5TXT")
	if err != nil {
		t.Fatal(err)
	}
	offer := xdr.LedgerKeyOffer{
		OwnerId: accountID,
		OfferId: 1,
	}

	Convey("Process ledger entry", t, func() {
		change := xdr.LedgerEntryChange{
			Type: xdr.LedgerEntryChangeTypeUpdated,
		}

		Convey("Wrong LedgerEntryChange", func() {
			s.processLedgerEntry(ctx, change, saleID)
			userConnector.AssertExpectations(t)
		})
	})

	Convey("Unique token", t, func() {
		emailAddress := "test_mail@gmail.com"
		var offerID uint64 = 1
		var saleID uint64 = 2
		expectedUniqueToken := fmt.Sprintf("%s:%d:%d:%s", emailAddress, offerID, saleID, s.config.RequestTokenSuffix)
		Convey("Expected unique token returns", func() {
			actualUniqueToken := s.buildUniqueToken(emailAddress, offerID, saleID)
			So(expectedUniqueToken, ShouldEqual, actualUniqueToken)
		})
		Convey("Request token suffix has changed", func() {
			s.config.RequestTokenSuffix += "s"
			actualUniqueToken := s.buildUniqueToken(emailAddress, offerID, saleID)
			So(expectedUniqueToken, ShouldNotEqual, actualUniqueToken)
		})
	})

	Convey("Email message build", t, func() {
		saleName := "Test sale"
		Convey("Template not found", func() {
			s.config.TemplateName += "s"
			_, err := s.buildEmailMessage(saleName)
			So(err, ShouldNotBeNil)
		})
		Convey("Successful build", func() {
			s.config.TemplateName = "io_cancelled.html"
			msg, err := s.buildEmailMessage(saleName)
			So(err, ShouldBeNil)
			So(msg, ShouldContainSubstring, config.TemplateLinkURL)
			So(msg, ShouldContainSubstring, saleName)
		})
	})

	Convey("User doesn't exist", t, func() {
		defer userConnector.AssertExpectations(t)
		defer saleConnector.AssertExpectations(t)

		userConnector.On("User", offer.OwnerId.Address()).Return(nil, nil).Once()
		emailUnit, err := s.processCancelledOrder(ctx, offer, saleID)

		So(emailUnit, ShouldBeNil)
		So(err, ShouldBeNil)
	})
	Convey("User exists", t, func() {
		user := &horizon.User{
			Type: "General",
			ID:   "User id",
			Attributes: horizon.UserAttributes{
				Email: "test_mail@gmail.com",
			},
		}

		userConnector.On("User", offer.OwnerId.Address()).Return(user, nil).Once()

		Convey("Sale doesn't exist", func() {
			saleConnector.On("SaleByID", saleID).Return(nil, nil).Once()

			emailUnit, err := s.processCancelledOrder(ctx, offer, saleID)
			So(emailUnit, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Sale exist", func() {
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

			emailUnit, err := s.processCancelledOrder(ctx, offer, saleID)
			So(err, ShouldBeNil)
			So(emailUnit.UniqueToken, ShouldEqual, s.buildUniqueToken(user.Attributes.Email, uint64(offer.OfferId), saleID))
			So(emailUnit.Payload.Destination, ShouldEqual, user.Attributes.Email)
			So(emailUnit.Payload.Subject, ShouldEqual, s.config.Subject)
			So(emailUnit.Payload.Message, ShouldEqual, msg)
		})
	})
}
