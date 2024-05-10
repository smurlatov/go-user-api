// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	models "user-api-service/internals/models"

	mock "github.com/stretchr/testify/mock"
)

// UserSaver is an autogenerated mock type for the UserSaver type
type UserSaver struct {
	mock.Mock
}

// SaveUser provides a mock function with given fields: user
func (_m *UserSaver) SaveUser(user models.User) (string, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for SaveUser")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(models.User) (string, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(models.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserSaver creates a new instance of UserSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserSaver {
	mock := &UserSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}