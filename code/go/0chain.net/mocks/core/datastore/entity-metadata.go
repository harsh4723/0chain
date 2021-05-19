// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	datastore "0chain.net/core/datastore"
	mock "github.com/stretchr/testify/mock"
)

// EntityMetadata is an autogenerated mock type for the EntityMetadata type
type EntityMetadata struct {
	mock.Mock
}

// GetDB provides a mock function with given fields:
func (_m *EntityMetadata) GetDB() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetIDColumnName provides a mock function with given fields:
func (_m *EntityMetadata) GetIDColumnName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetName provides a mock function with given fields:
func (_m *EntityMetadata) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetStore provides a mock function with given fields:
func (_m *EntityMetadata) GetStore() datastore.Store {
	ret := _m.Called()

	var r0 datastore.Store
	if rf, ok := ret.Get(0).(func() datastore.Store); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(datastore.Store)
		}
	}

	return r0
}

// Instance provides a mock function with given fields:
func (_m *EntityMetadata) Instance() datastore.Entity {
	ret := _m.Called()

	var r0 datastore.Entity
	if rf, ok := ret.Get(0).(func() datastore.Entity); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(datastore.Entity)
		}
	}

	return r0
}