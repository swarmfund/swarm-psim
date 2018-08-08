// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"
import regources "gitlab.com/tokend/regources"

// Infoer is an autogenerated mock type for the Infoer type
type Infoer struct {
	mock.Mock
}

// Info provides a mock function with given fields:
func (_m *Infoer) Info() (*regources.Info, error) {
	ret := _m.Called()

	var r0 *regources.Info
	if rf, ok := ret.Get(0).(func() *regources.Info); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*regources.Info)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
