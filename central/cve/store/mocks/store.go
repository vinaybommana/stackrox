// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// GetAll mocks base method
func (m *MockStore) GetAll() ([]*storage.CVE, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*storage.CVE)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockStoreMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStore)(nil).GetAll))
}

// Count mocks base method
func (m *MockStore) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockStoreMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockStore)(nil).Count))
}

// Get mocks base method
func (m *MockStore) Get(id string) (*storage.CVE, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*storage.CVE)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get
func (mr *MockStoreMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), id)
}

// GetBatch mocks base method
func (m *MockStore) GetBatch(ids []string) ([]*storage.CVE, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatch", ids)
	ret0, _ := ret[0].([]*storage.CVE)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBatch indicates an expected call of GetBatch
func (mr *MockStoreMockRecorder) GetBatch(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatch", reflect.TypeOf((*MockStore)(nil).GetBatch), ids)
}

// Exists mocks base method
func (m *MockStore) Exists(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists
func (mr *MockStoreMockRecorder) Exists(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockStore)(nil).Exists), id)
}

// Upsert mocks base method
func (m *MockStore) Upsert(cves ...*storage.CVE) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range cves {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Upsert", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert
func (mr *MockStoreMockRecorder) Upsert(cves ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockStore)(nil).Upsert), cves...)
}

// Delete mocks base method
func (m *MockStore) Delete(ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockStoreMockRecorder) Delete(ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStore)(nil).Delete), ids...)
}

// GetTxnCount mocks base method
func (m *MockStore) GetTxnCount() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxnCount")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxnCount indicates an expected call of GetTxnCount
func (mr *MockStoreMockRecorder) GetTxnCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxnCount", reflect.TypeOf((*MockStore)(nil).GetTxnCount))
}

// IncTxnCount mocks base method
func (m *MockStore) IncTxnCount() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncTxnCount")
	ret0, _ := ret[0].(error)
	return ret0
}

// IncTxnCount indicates an expected call of IncTxnCount
func (mr *MockStoreMockRecorder) IncTxnCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncTxnCount", reflect.TypeOf((*MockStore)(nil).IncTxnCount))
}
