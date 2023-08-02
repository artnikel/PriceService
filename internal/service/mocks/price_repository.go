// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/artnikel/PriceService/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// PriceRepository is an autogenerated mock type for the PriceRepository type
type PriceRepository struct {
	mock.Mock
}

// ReadPrices provides a mock function with given fields: ctx
func (_m *PriceRepository) ReadPrices(ctx context.Context) ([]*model.Action, error) {
	ret := _m.Called(ctx)

	var r0 []*model.Action
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Action); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Action)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPriceRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPriceRepository creates a new instance of PriceRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPriceRepository(t mockConstructorTestingTNewPriceRepository) *PriceRepository {
	mock := &PriceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}