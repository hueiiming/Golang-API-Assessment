// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "Golang-API-Assessment/pkg/types"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// ClearTables provides a mock function with given fields:
func (_m *Repository) ClearTables() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ClearTables")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCommonStudents provides a mock function with given fields: teachers
func (_m *Repository) GetCommonStudents(teachers []string) ([]string, error) {
	ret := _m.Called(teachers)

	if len(ret) == 0 {
		panic("no return value specified for GetCommonStudents")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) ([]string, error)); ok {
		return rf(teachers)
	}
	if rf, ok := ret.Get(0).(func([]string) []string); ok {
		r0 = rf(teachers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(teachers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotification provides a mock function with given fields: request
func (_m *Repository) GetNotification(request *types.NotificationRequest) ([]string, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for GetNotification")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(*types.NotificationRequest) ([]string, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*types.NotificationRequest) []string); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(*types.NotificationRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStudentID provides a mock function with given fields: studentEmail
func (_m *Repository) GetStudentID(studentEmail string) (int, error) {
	ret := _m.Called(studentEmail)

	if len(ret) == 0 {
		panic("no return value specified for GetStudentID")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int, error)); ok {
		return rf(studentEmail)
	}
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(studentEmail)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(studentEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTeacherID provides a mock function with given fields: teacherEmail
func (_m *Repository) GetTeacherID(teacherEmail string) (int, error) {
	ret := _m.Called(teacherEmail)

	if len(ret) == 0 {
		panic("no return value specified for GetTeacherID")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int, error)); ok {
		return rf(teacherEmail)
	}
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(teacherEmail)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(teacherEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PopulateTables provides a mock function with given fields: teacherEmails, studentEmails
func (_m *Repository) PopulateTables(teacherEmails []string, studentEmails []string) error {
	ret := _m.Called(teacherEmails, studentEmails)

	if len(ret) == 0 {
		panic("no return value specified for PopulateTables")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, []string) error); ok {
		r0 = rf(teacherEmails, studentEmails)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Registration provides a mock function with given fields: teacherID, studentID
func (_m *Repository) Registration(teacherID int, studentID []int) error {
	ret := _m.Called(teacherID, studentID)

	if len(ret) == 0 {
		panic("no return value specified for Registration")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []int) error); ok {
		r0 = rf(teacherID, studentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Suspension provides a mock function with given fields: studentID
func (_m *Repository) Suspension(studentID int) error {
	ret := _m.Called(studentID)

	if len(ret) == 0 {
		panic("no return value specified for Suspension")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(studentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
