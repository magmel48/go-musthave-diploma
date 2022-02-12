// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	users "github.com/magmel48/go-musthave-diploma/internal/users"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *Repository) Create(ctx context.Context, user users.User) (int64, error) {
	ret := _m.Called(ctx, user)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, users.User) int64); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, login
func (_m *Repository) Find(ctx context.Context, login string) (*users.User, error) {
	ret := _m.Called(ctx, login)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *users.User); ok {
		r0 = rf(ctx, login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
