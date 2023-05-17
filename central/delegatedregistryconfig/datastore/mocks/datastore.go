// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
)

// MockDataStore is a mock of DataStore interface.
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore.
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance.
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// GetConfig mocks base method.
func (m *MockDataStore) GetConfig(arg0 context.Context) (*storage.DelegatedRegistryConfig, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig", arg0)
	ret0, _ := ret[0].(*storage.DelegatedRegistryConfig)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockDataStoreMockRecorder) GetConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockDataStore)(nil).GetConfig), arg0)
}

// UpsertConfig mocks base method.
func (m *MockDataStore) UpsertConfig(arg0 context.Context, arg1 *storage.DelegatedRegistryConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertConfig indicates an expected call of UpsertConfig.
func (mr *MockDataStoreMockRecorder) UpsertConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertConfig", reflect.TypeOf((*MockDataStore)(nil).UpsertConfig), arg0, arg1)
}
