package bearer

import (
	"context"

	"gitlab.com/swarmfund/horizon-connector/v2"
	"github.com/stretchr/testify/mock"
)

type SaleStateCheckerMock struct {
	mock.Mock
}

func (sm *SaleStateCheckerMock) GetSales() ([]horizon.Sale, error) {
	a := sm.Called()
	sales, _ := a.Get(0).([]horizon.Sale)
	return sales, a.Error(1)
}

func (sm *SaleStateCheckerMock) GetHorizonInfo() (info *horizon.Info, err error) {
	a := sm.Called()
	horizonInfo, _ := a.Get(0).(*horizon.Info)
	return horizonInfo, a.Error(1)
}

func (sm *SaleStateCheckerMock) BuildTx(info *horizon.Info, saleID uint64) (string, error) {
	args := sm.Called(info, saleID)
	result, _ := args.Get(0).(string)
	return result, args.Error(1)
}

func (sm *SaleStateCheckerMock) SubmitTx(ctx context.Context, envelope string) horizon.SubmitResult {
	args := sm.Called(ctx, envelope)
	result := args.Get(0).(horizon.SubmitResult)
	return result
}