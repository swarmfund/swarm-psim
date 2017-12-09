package txhandler

import (
	"math/rand"
	"testing"

	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/psim/psim/taxman/internal/txhandler/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func TestHandler(t *testing.T) {
	Convey("Given valid handler", t, func() {
		statable := &mocks.Statable{}
		defer statable.AssertExpectations(t)
		handler := NewHandler(statable, map[string]bool{}, logan.New())
		Convey("Should skip tx", func() {
			tx := horizon.Transaction{
				ID: "should_skip",
			}
			handler.txsToSkip[tx.ID] = true
			err := handler.Handle(tx)
			So(err, ShouldBeNil)
		})
		Convey("One of the handlers returns error", func() {
			tx := horizon.Transaction{}
			txHandlerMock := &mocks.HorizonTxHandler{}
			txHandlerMock.On("Handle", tx).Return(errors.New("failed to handle tx")).Once()
			handler.handlers = []HorizonTxHandler{txHandlerMock}
			err := handler.Handle(tx)
			So(err, ShouldNotBeNil)
		})
		Convey("Success", func() {
			tx := horizon.Transaction{
				ID:     "should_not_skip",
				Ledger: rand.Int63(),
			}
			statable.On("SetLedger", tx.Ledger).Once()
			txHandlerMock := &mocks.HorizonTxHandler{}
			txHandlerMock.On("Handle", tx).Return(nil).Once()
			handler.handlers = []HorizonTxHandler{txHandlerMock}
			err := handler.Handle(tx)
			So(err, ShouldBeNil)
		})
	})
}
