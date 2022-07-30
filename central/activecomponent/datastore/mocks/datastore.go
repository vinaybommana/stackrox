// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
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

// DeleteBatch mocks base method.
func (m *MockDataStore) DeleteBatch(ctx context.Context, ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteBatch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch.
func (mr *MockDataStoreMockRecorder) DeleteBatch(ctx interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockDataStore)(nil).DeleteBatch), varargs...)
}

// Exists mocks base method.
func (m *MockDataStore) Exists(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockDataStoreMockRecorder) Exists(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockDataStore)(nil).Exists), ctx, id)
}

// Get mocks base method.
func (m *MockDataStore) Get(ctx context.Context, id string) (*storage.ActiveComponent, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*storage.ActiveComponent)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockDataStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDataStore)(nil).Get), ctx, id)
}

// GetBatch mocks base method.
func (m *MockDataStore) GetBatch(ctx context.Context, ids []string) ([]*storage.ActiveComponent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatch", ctx, ids)
	ret0, _ := ret[0].([]*storage.ActiveComponent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatch indicates an expected call of GetBatch.
func (mr *MockDataStoreMockRecorder) GetBatch(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatch", reflect.TypeOf((*MockDataStore)(nil).GetBatch), ctx, ids)
}

// Search mocks base method.
func (m *MockDataStore) Search(ctx context.Context, query *auxpb.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockDataStoreMockRecorder) Search(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, query)
}

// SearchRawActiveComponents mocks base method.
func (m *MockDataStore) SearchRawActiveComponents(ctx context.Context, q *auxpb.Query) ([]*storage.ActiveComponent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawActiveComponents", ctx, q)
	ret0, _ := ret[0].([]*storage.ActiveComponent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawActiveComponents indicates an expected call of SearchRawActiveComponents.
func (mr *MockDataStoreMockRecorder) SearchRawActiveComponents(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawActiveComponents", reflect.TypeOf((*MockDataStore)(nil).SearchRawActiveComponents), ctx, q)
}

// UpsertBatch mocks base method.
func (m *MockDataStore) UpsertBatch(ctx context.Context, activeComponents []*storage.ActiveComponent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertBatch", ctx, activeComponents)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertBatch indicates an expected call of UpsertBatch.
func (mr *MockDataStoreMockRecorder) UpsertBatch(ctx, activeComponents interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertBatch", reflect.TypeOf((*MockDataStore)(nil).UpsertBatch), ctx, activeComponents)
}
