package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/tokend/horizon-connector"
)

// TransactionConnector is an autogenerated mock type for the TransactionConnector type
type TransactionConnector struct {
	mock.Mock
}

// TransactionByID provides a mock function with given fields: txID
func (_m *TransactionConnector) TransactionByID(txID string) (*horizon.Transaction, error) {
	ret := _m.Called(txID)

	var r0 *horizon.Transaction
	if rf, ok := ret.Get(0).(func(string) *horizon.Transaction); ok {
		r0 = rf(txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*horizon.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
