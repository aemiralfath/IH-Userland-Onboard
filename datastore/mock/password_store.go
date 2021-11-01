// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aemiralfath/IH-Userland-Onboard/datastore (interfaces: PasswordStore)

// Package mock_datastore is a generated GoMock package.
package mock_datastore

import (
	context "context"
	datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPasswordStore is a mock of PasswordStore interface
type MockPasswordStore struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordStoreMockRecorder
}

// MockPasswordStoreMockRecorder is the mock recorder for MockPasswordStore
type MockPasswordStoreMockRecorder struct {
	mock *MockPasswordStore
}

// NewMockPasswordStore creates a new mock instance
func NewMockPasswordStore(ctrl *gomock.Controller) *MockPasswordStore {
	mock := &MockPasswordStore{ctrl: ctrl}
	mock.recorder = &MockPasswordStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPasswordStore) EXPECT() *MockPasswordStoreMockRecorder {
	return m.recorder
}

// AddNewPassword mocks base method
func (m *MockPasswordStore) AddNewPassword(arg0 context.Context, arg1 *datastore.Password, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewPassword indicates an expected call of AddNewPassword
func (mr *MockPasswordStoreMockRecorder) AddNewPassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewPassword", reflect.TypeOf((*MockPasswordStore)(nil).AddNewPassword), arg0, arg1, arg2)
}

// GetLastThreePassword mocks base method
func (m *MockPasswordStore) GetLastThreePassword(arg0 context.Context, arg1 float64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastThreePassword", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastThreePassword indicates an expected call of GetLastThreePassword
func (mr *MockPasswordStoreMockRecorder) GetLastThreePassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastThreePassword", reflect.TypeOf((*MockPasswordStore)(nil).GetLastThreePassword), arg0, arg1)
}
