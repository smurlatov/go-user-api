// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	models "user-api-service/internals/models"

	mock "github.com/stretchr/testify/mock"
)

// UserUpdater is an autogenerated mock type for the UserUpdater type
type UserUpdater struct {
	mock.Mock
}

// UpdateUser provides a mock function with given fields: user, id
func (_m *UserUpdater) UpdateUser(user models.User, id string) error {
	ret := _m.Called(user, id)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(models.User, string) error); ok {
		r0 = rf(user, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserUpdater creates a new instance of UserUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserUpdater {
	mock := &UserUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}