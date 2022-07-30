// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
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

// Count mocks base method.
func (m *MockDataStore) Count(ctx context.Context, q *auxpb.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, q)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockDataStoreMockRecorder) Count(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDataStore)(nil).Count), ctx, q)
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
func (m *MockDataStore) Get(ctx context.Context, id string) (*storage.ComponentCVEEdge, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*storage.ComponentCVEEdge)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockDataStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDataStore)(nil).Get), ctx, id)
}

// Search mocks base method.
func (m *MockDataStore) Search(ctx context.Context, q *auxpb.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockDataStoreMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// SearchEdges mocks base method.
func (m *MockDataStore) SearchEdges(ctx context.Context, q *auxpb.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchEdges", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchEdges indicates an expected call of SearchEdges.
func (mr *MockDataStoreMockRecorder) SearchEdges(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchEdges", reflect.TypeOf((*MockDataStore)(nil).SearchEdges), ctx, q)
}

// SearchRawEdges mocks base method.
func (m *MockDataStore) SearchRawEdges(ctx context.Context, q *auxpb.Query) ([]*storage.ComponentCVEEdge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawEdges", ctx, q)
	ret0, _ := ret[0].([]*storage.ComponentCVEEdge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawEdges indicates an expected call of SearchRawEdges.
func (mr *MockDataStoreMockRecorder) SearchRawEdges(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawEdges", reflect.TypeOf((*MockDataStore)(nil).SearchRawEdges), ctx, q)
}
