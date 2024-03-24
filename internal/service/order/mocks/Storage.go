// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/DivPro/app/pkg/entity"
	mock "github.com/stretchr/testify/mock"

	storage "github.com/DivPro/app/internal/storage"

	time "time"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// CreateOrder provides a mock function with given fields: ctx, _a1
func (_m *Storage) CreateOrder(ctx context.Context, _a1 *entity.Order) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Order) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAvailability provides a mock function with given fields: ctx, date, hotelID, roomID
func (_m *Storage) FindAvailability(ctx context.Context, date time.Time, hotelID string, roomID string) (*entity.RoomAvailability, error) {
	ret := _m.Called(ctx, date, hotelID, roomID)

	if len(ret) == 0 {
		panic("no return value specified for FindAvailability")
	}

	var r0 *entity.RoomAvailability
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, string, string) (*entity.RoomAvailability, error)); ok {
		return rf(ctx, date, hotelID, roomID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, string, string) *entity.RoomAvailability); ok {
		r0 = rf(ctx, date, hotelID, roomID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RoomAvailability)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time, string, string) error); ok {
		r1 = rf(ctx, date, hotelID, roomID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Tx provides a mock function with given fields: ctx, fn
func (_m *Storage) Tx(ctx context.Context, fn storage.TxFn) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for Tx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, storage.TxFn) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAvailabilities provides a mock function with given fields: ctx, availabilities
func (_m *Storage) UpdateAvailabilities(ctx context.Context, availabilities []*entity.RoomAvailability) error {
	ret := _m.Called(ctx, availabilities)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAvailabilities")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*entity.RoomAvailability) error); ok {
		r0 = rf(ctx, availabilities)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
