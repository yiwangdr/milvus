// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// TxnKV is an autogenerated mock type for the TxnKV type
type TxnKV struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *TxnKV) Close() {
	_m.Called()
}

// Load provides a mock function with given fields: key
func (_m *TxnKV) Load(key string) (string, error) {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadWithPrefix provides a mock function with given fields: key
func (_m *TxnKV) LoadWithPrefix(key string) ([]string, []string, error) {
	ret := _m.Called(key)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 []string
	if rf, ok := ret.Get(1).(func(string) []string); ok {
		r1 = rf(key)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(key)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MultiLoad provides a mock function with given fields: keys
func (_m *TxnKV) MultiLoad(keys []string) ([]string, error) {
	ret := _m.Called(keys)

	var r0 []string
	if rf, ok := ret.Get(0).(func([]string) []string); ok {
		r0 = rf(keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultiRemove provides a mock function with given fields: keys
func (_m *TxnKV) MultiRemove(keys []string) error {
	ret := _m.Called(keys)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(keys)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MultiRemoveWithPrefix provides a mock function with given fields: keys
func (_m *TxnKV) MultiRemoveWithPrefix(keys []string) error {
	ret := _m.Called(keys)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(keys)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MultiSave provides a mock function with given fields: kvs
func (_m *TxnKV) MultiSave(kvs map[string]string) error {
	ret := _m.Called(kvs)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]string) error); ok {
		r0 = rf(kvs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MultiSaveAndRemove provides a mock function with given fields: saves, removals
func (_m *TxnKV) MultiSaveAndRemove(saves map[string]string, removals []string) error {
	ret := _m.Called(saves, removals)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]string, []string) error); ok {
		r0 = rf(saves, removals)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MultiSaveAndRemoveWithPrefix provides a mock function with given fields: saves, removals
func (_m *TxnKV) MultiSaveAndRemoveWithPrefix(saves map[string]string, removals []string) error {
	ret := _m.Called(saves, removals)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]string, []string) error); ok {
		r0 = rf(saves, removals)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: key
func (_m *TxnKV) Remove(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveWithPrefix provides a mock function with given fields: key
func (_m *TxnKV) RemoveWithPrefix(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: key, value
func (_m *TxnKV) Save(key string, value string) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTxnKV interface {
	mock.TestingT
	Cleanup(func())
}

// NewTxnKV creates a new instance of TxnKV. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTxnKV(t mockConstructorTestingTNewTxnKV) *TxnKV {
	mock := &TxnKV{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
