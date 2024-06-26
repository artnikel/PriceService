// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	proto "github.com/artnikel/PriceService/proto"
)

// PriceServiceClient is an autogenerated mock type for the PriceServiceClient type
type PriceServiceClient struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: ctx, in, opts
func (_m *PriceServiceClient) Subscribe(ctx context.Context, in *proto.SubscribeRequest, opts ...grpc.CallOption) (proto.PriceService_SubscribeClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 proto.PriceService_SubscribeClient
	if rf, ok := ret.Get(0).(func(context.Context, *proto.SubscribeRequest, ...grpc.CallOption) proto.PriceService_SubscribeClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(proto.PriceService_SubscribeClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.SubscribeRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPriceServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewPriceServiceClient creates a new instance of PriceServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPriceServiceClient(t mockConstructorTestingTNewPriceServiceClient) *PriceServiceClient {
	mock := &PriceServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
