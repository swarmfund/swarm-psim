// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"
import s3 "github.com/aws/aws-sdk-go/service/s3"
import s3manager "github.com/aws/aws-sdk-go/service/s3/s3manager"

// TemplateDownloader is an autogenerated mock type for the TemplateDownloader type
type TemplateDownloader struct {
	mock.Mock
}

// Download provides a mock function with given fields: w, input, options
func (_m *TemplateDownloader) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, w, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int64
	if rf, ok := ret.Get(0).(func(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) int64); ok {
		r0 = rf(w, input, options...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) error); ok {
		r1 = rf(w, input, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
