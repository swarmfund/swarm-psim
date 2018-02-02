package bearer

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

func TestCheckSalesStateHelper(t *testing.T) {
	checkSalesStateHelper := CheckSalesStateHelperMock{}
	ctx := context.Background()
	log := logan.New()
	config := Config{}
	service := New(config, log, &checkSalesStateHelper)

	Convey("Try to check sales", t, func() {
		Convey("Failed to get sales", func() {
			expectedError := errors.New("failed to get sales")
			checkSalesStateHelper.On("GetSales").Return(nil, expectedError).Once()

			err := service.checkSalesState(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, expectedError, errors.Cause(err))
		})
		Convey("Got sales successfully", func() {
			sale := horizon.Sale{ID: 1}
			var sales []horizon.Sale
			sales = append(sales, sale)
			checkSalesStateHelper.On("GetSales").Return(sales, nil).Once()

			Convey("Failed to get horizon info", func() {
				expectedError := errors.New("failed to get horizon info")
				checkSalesStateHelper.On("GetHorizonInfo").Return(nil, expectedError).Once()

				err := service.checkSalesState(ctx)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError, errors.Cause(err))
			})
			Convey("Got horizon info successfully", func() {
				info := horizon.Info{
					Passphrase:         "Caput draconis",
					MasterAccountID:    "MasterAccountID",
					TXExpirationPeriod: 42,
				}
				checkSalesStateHelper.On("GetHorizonInfo").Return(&info, nil).Once()

				Convey("Failed to build tx", func() {
					expectedError := errors.New("failed to build tx")
					checkSalesStateHelper.On("BuildTx", &info, sale.ID).Return("", expectedError).Once()

					err := service.checkSalesState(ctx)
					assert.NotNil(t, err)
					assert.Equal(t, expectedError, errors.Cause(err))
				})
				Convey("Tx built successfully", func() {
					envelope := "Tx build success"
					checkSalesStateHelper.On("BuildTx", &info, sale.ID).Return(envelope, nil).Once()

					Convey("Failed to submit tx", func() {
						expectedError := errors.New("failed to submit tx")
						checkSalesStateHelper.On("SubmitTx", ctx, envelope).Return(horizon.SubmitResult{Err: expectedError}).Once()

						err := service.checkSalesState(ctx)
						assert.NotNil(t, err)
						assert.Equal(t, expectedError, errors.Cause(err))
					})
					Convey("CheckSaleState success", func() {
						checkSalesStateHelper.On("SubmitTx", ctx, envelope).Return(horizon.SubmitResult{Err: nil}).Once()

						err := service.checkSalesState(ctx)
						assert.Nil(t, err)
					})
				})
			})
		})
	})
}
