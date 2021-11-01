// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aemiralfath/IH-Userland-Onboard/datastore (interfaces: ProfileStore)

// Package mock_datastore is a generated GoMock package.
package mock_datastore

import (
	context "context"
	datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockProfileStore is a mock of ProfileStore interface
type MockProfileStore struct {
	ctrl     *gomock.Controller
	recorder *MockProfileStoreMockRecorder
}

// MockProfileStoreMockRecorder is the mock recorder for MockProfileStore
type MockProfileStoreMockRecorder struct {
	mock *MockProfileStore
}

// NewMockProfileStore creates a new mock instance
func NewMockProfileStore(ctrl *gomock.Controller) *MockProfileStore {
	mock := &MockProfileStore{ctrl: ctrl}
	mock.recorder = &MockProfileStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProfileStore) EXPECT() *MockProfileStoreMockRecorder {
	return m.recorder
}

// AddNewProfile mocks base method
func (m *MockProfileStore) AddNewProfile(arg0 context.Context, arg1 *datastore.Profile, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewProfile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewProfile indicates an expected call of AddNewProfile
func (mr *MockProfileStoreMockRecorder) AddNewProfile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewProfile", reflect.TypeOf((*MockProfileStore)(nil).AddNewProfile), arg0, arg1, arg2)
}

// GetProfile mocks base method
func (m *MockProfileStore) GetProfile(arg0 context.Context, arg1 float64) (*datastore.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", arg0, arg1)
	ret0, _ := ret[0].(*datastore.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile
func (mr *MockProfileStoreMockRecorder) GetProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockProfileStore)(nil).GetProfile), arg0, arg1)
}

// UpdatePicture mocks base method
func (m *MockProfileStore) UpdatePicture(arg0 context.Context, arg1 *datastore.Profile, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePicture", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePicture indicates an expected call of UpdatePicture
func (mr *MockProfileStoreMockRecorder) UpdatePicture(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePicture", reflect.TypeOf((*MockProfileStore)(nil).UpdatePicture), arg0, arg1, arg2)
}

// UpdateProfile mocks base method
func (m *MockProfileStore) UpdateProfile(arg0 context.Context, arg1 *datastore.Profile, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfile indicates an expected call of UpdateProfile
func (mr *MockProfileStoreMockRecorder) UpdateProfile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockProfileStore)(nil).UpdateProfile), arg0, arg1, arg2)
}
