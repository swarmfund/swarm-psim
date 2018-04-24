package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/tokend/horizon-connector"
)

// SaleConnector is an autogenerated mock type for the SaleConnector type
type SaleConnector struct {
	mock.Mock
}

// SaleByID provides a mock function with given fields: saleID
func (_m *SaleConnector) SaleByID(saleID uint64) (*horizon.Sale, error) {
	ret := _m.Called(saleID)

	var r0 *horizon.Sale
	if rf, ok := ret.Get(0).(func(uint64) *horizon.Sale); ok {
		r0 = rf(saleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*horizon.Sale)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(saleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
