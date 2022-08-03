// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
)

// MockAttackReadOnlyDataStore is a mock of AttackReadOnlyDataStore interface.
type MockAttackReadOnlyDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockAttackReadOnlyDataStoreMockRecorder
}

// MockAttackReadOnlyDataStoreMockRecorder is the mock recorder for MockAttackReadOnlyDataStore.
type MockAttackReadOnlyDataStoreMockRecorder struct {
	mock *MockAttackReadOnlyDataStore
}

// NewMockAttackReadOnlyDataStore creates a new mock instance.
func NewMockAttackReadOnlyDataStore(ctrl *gomock.Controller) *MockAttackReadOnlyDataStore {
	mock := &MockAttackReadOnlyDataStore{ctrl: ctrl}
	mock.recorder = &MockAttackReadOnlyDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAttackReadOnlyDataStore) EXPECT() *MockAttackReadOnlyDataStoreMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockAttackReadOnlyDataStore) Get(id string) (*storage.MitreAttackVector, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*storage.MitreAttackVector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAttackReadOnlyDataStoreMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAttackReadOnlyDataStore)(nil).Get), id)
}

// GetAll mocks base method.
func (m *MockAttackReadOnlyDataStore) GetAll() []*storage.MitreAttackVector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*storage.MitreAttackVector)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockAttackReadOnlyDataStoreMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockAttackReadOnlyDataStore)(nil).GetAll))
}