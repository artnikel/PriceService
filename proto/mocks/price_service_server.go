// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	proto "github.com/artnikel/PriceService/proto"
	mock "github.com/stretchr/testify/mock"
)

// PriceServiceServer is an autogenerated mock type for the PriceServiceServer type
type PriceServiceServer struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: _a0, _a1
func (_m *PriceServiceServer) Subscribe(_a0 *proto.SubscribeRequest, _a1 proto.PriceService_SubscribeServer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*proto.SubscribeRequest, proto.PriceService_SubscribeServer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mustEmbedUnimplementedPriceServiceServer provides a mock function with given fields:
func (_m *PriceServiceServer) mustEmbedUnimplementedPriceServiceServer() {
	_m.Called()
}

type mockConstructorTestingTNewPriceServiceServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewPriceServiceServer creates a new instance of PriceServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPriceServiceServer(t mockConstructorTestingTNewPriceServiceServer) *PriceServiceServer {
	mock := &PriceServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
