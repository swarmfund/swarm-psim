package bearer

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

func TestSalesStateChecker(t *testing.T) {
	salesStateChecker := SaleStateCheckerMock{}
	ctx := context.Background()
	log := logan.New()
	config := Config{}
	service := New(config, log, &salesStateChecker)

	Convey("Try to check sales", t, func() {
		Convey("Failed to get sales", func() {
			expectedError := errors.New("failed to get sales")
			salesStateChecker.On("GetSales").Return(nil, expectedError).Once()
			err := service.checkSaleState(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, expectedError, errors.Cause(err))
		})
		Convey("Failed to get horizon info", func() {
			expectedError := errors.New("failed to get horizon info")
			sale := horizon.Sale{1}
			var sales []horizon.Sale
			sales = append(sales, sale)
			salesStateChecker.On("GetSales").Return(sales, nil).Once()
			salesStateChecker.On("GetHorizonInfo").Return(nil, expectedError).Once()
			err := service.checkSaleState(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, expectedError, errors.Cause(err))
		})
		Convey("Failed to marshal tx", func() {
			expectedError := errors.New("failed to marshal tx")
			sale := horizon.Sale{1}
			var sales []horizon.Sale
			sales = append(sales, sale)
			salesStateChecker.On("GetSales").Return(sales, nil).Once()
			info := horizon.Info{"Caput draconis", "MasterAccountID", 42}
			infoPtr := &info
			salesStateChecker.On("GetHorizonInfo").Return(infoPtr, nil).Once()
			salesStateChecker.On("BuildTx", infoPtr, sale.ID).Return("", expectedError).Once()
			err := service.checkSaleState(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, expectedError, errors.Cause(err))
		})
		Convey("Failed to submit tx", func() {
			expectedError := errors.New("failed to submit tx")
			sale := horizon.Sale{1}
			var sales []horizon.Sale
			sales = append(sales, sale)
			salesStateChecker.On("GetSales").Return(sales, nil).Once()
			info := horizon.Info{"Caput draconis", "MasterAccountID", 42}
			infoPtr := &info
			salesStateChecker.On("GetHorizonInfo").Return(infoPtr, nil).Once()
			envelope := "Tx build success"
			salesStateChecker.On("BuildTx", infoPtr, sale.ID).Return(envelope, nil).Once()
			salesStateChecker.On("SubmitTx", ctx, envelope).Return(horizon.SubmitResult{Err: expectedError}).Once()
			err := service.checkSaleState(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, expectedError, errors.Cause(err))
		})
		Convey("Success", func() {
			sale := horizon.Sale{1}
			var sales []horizon.Sale
			sales = append(sales, sale)
			salesStateChecker.On("GetSales").Return(sales, nil).Once()
			info := horizon.Info{"Caput draconis", "MasterAccountID", 42}
			infoPtr := &info
			salesStateChecker.On("GetHorizonInfo").Return(infoPtr, nil).Once()
			envelope := "Tx build success"
			salesStateChecker.On("BuildTx", infoPtr, sale.ID).Return(envelope, nil).Once()
			var bytes []byte
			var opCodes []string
			submitResult := horizon.SubmitResult{
				Err:         nil,
				RawResponse: bytes,
				TXCode:      "txSUCCESS",
				OpCodes:     opCodes,
			}
			salesStateChecker.On("SubmitTx", ctx, envelope).Return(submitResult).Once()
			err := service.checkSaleState(ctx)
			assert.Nil(t, err)
		})
	})
}
