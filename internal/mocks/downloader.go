// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "httpfs/internal/entities"

	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Downloader is an autogenerated mock type for the Downloader type
type Downloader struct {
	mock.Mock
}

type Downloader_Expecter struct {
	mock *mock.Mock
}

func (_m *Downloader) EXPECT() *Downloader_Expecter {
	return &Downloader_Expecter{mock: &_m.Mock}
}

// Download provides a mock function with given fields: _a0, _a1, _a2
func (_m *Downloader) Download(_a0 context.Context, _a1 entities.Hash, _a2 io.Writer) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Download")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entities.Hash, io.Writer) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Downloader_Download_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Download'
type Downloader_Download_Call struct {
	*mock.Call
}

// Download is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 entities.Hash
//   - _a2 io.Writer
func (_e *Downloader_Expecter) Download(_a0 interface{}, _a1 interface{}, _a2 interface{}) *Downloader_Download_Call {
	return &Downloader_Download_Call{Call: _e.mock.On("Download", _a0, _a1, _a2)}
}

func (_c *Downloader_Download_Call) Run(run func(_a0 context.Context, _a1 entities.Hash, _a2 io.Writer)) *Downloader_Download_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entities.Hash), args[2].(io.Writer))
	})
	return _c
}

func (_c *Downloader_Download_Call) Return(_a0 error) *Downloader_Download_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Downloader_Download_Call) RunAndReturn(run func(context.Context, entities.Hash, io.Writer) error) *Downloader_Download_Call {
	_c.Call.Return(run)
	return _c
}

// NewDownloader creates a new instance of Downloader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDownloader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Downloader {
	mock := &Downloader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
