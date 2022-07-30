// Code generated by MockGen. DO NOT EDIT.
// Source: indexer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

// MockIndexer is a mock of Indexer interface.
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer.
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance.
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// AddPolicyCategories mocks base method.
func (m *MockIndexer) AddPolicyCategories(policycategories []*storage.PolicyCategory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPolicyCategories", policycategories)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPolicyCategories indicates an expected call of AddPolicyCategories.
func (mr *MockIndexerMockRecorder) AddPolicyCategories(policycategories interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPolicyCategories", reflect.TypeOf((*MockIndexer)(nil).AddPolicyCategories), policycategories)
}

// AddPolicyCategory mocks base method.
func (m *MockIndexer) AddPolicyCategory(policycategory *storage.PolicyCategory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPolicyCategory", policycategory)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPolicyCategory indicates an expected call of AddPolicyCategory.
func (mr *MockIndexerMockRecorder) AddPolicyCategory(policycategory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPolicyCategory", reflect.TypeOf((*MockIndexer)(nil).AddPolicyCategory), policycategory)
}

// Count mocks base method.
func (m *MockIndexer) Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{q}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Count", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIndexerMockRecorder) Count(q interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{q}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIndexer)(nil).Count), varargs...)
}

// DeletePolicyCategories mocks base method.
func (m *MockIndexer) DeletePolicyCategories(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePolicyCategories", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePolicyCategories indicates an expected call of DeletePolicyCategories.
func (mr *MockIndexerMockRecorder) DeletePolicyCategories(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePolicyCategories", reflect.TypeOf((*MockIndexer)(nil).DeletePolicyCategories), ids)
}

// DeletePolicyCategory mocks base method.
func (m *MockIndexer) DeletePolicyCategory(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePolicyCategory", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePolicyCategory indicates an expected call of DeletePolicyCategory.
func (mr *MockIndexerMockRecorder) DeletePolicyCategory(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePolicyCategory", reflect.TypeOf((*MockIndexer)(nil).DeletePolicyCategory), id)
}

// MarkInitialIndexingComplete mocks base method.
func (m *MockIndexer) MarkInitialIndexingComplete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkInitialIndexingComplete")
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkInitialIndexingComplete indicates an expected call of MarkInitialIndexingComplete.
func (mr *MockIndexerMockRecorder) MarkInitialIndexingComplete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkInitialIndexingComplete", reflect.TypeOf((*MockIndexer)(nil).MarkInitialIndexingComplete))
}

// NeedsInitialIndexing mocks base method.
func (m *MockIndexer) NeedsInitialIndexing() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NeedsInitialIndexing")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NeedsInitialIndexing indicates an expected call of NeedsInitialIndexing.
func (mr *MockIndexerMockRecorder) NeedsInitialIndexing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NeedsInitialIndexing", reflect.TypeOf((*MockIndexer)(nil).NeedsInitialIndexing))
}

// Search mocks base method.
func (m *MockIndexer) Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{q}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Search", varargs...)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockIndexerMockRecorder) Search(q interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{q}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIndexer)(nil).Search), varargs...)
}
