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

// AddImageCVEEdge mocks base method.
func (m *MockIndexer) AddImageCVEEdge(imagecveedge *storage.ImageCVEEdge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImageCVEEdge", imagecveedge)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImageCVEEdge indicates an expected call of AddImageCVEEdge.
func (mr *MockIndexerMockRecorder) AddImageCVEEdge(imagecveedge interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImageCVEEdge", reflect.TypeOf((*MockIndexer)(nil).AddImageCVEEdge), imagecveedge)
}

// AddImageCVEEdges mocks base method.
func (m *MockIndexer) AddImageCVEEdges(imagecveedges []*storage.ImageCVEEdge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImageCVEEdges", imagecveedges)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImageCVEEdges indicates an expected call of AddImageCVEEdges.
func (mr *MockIndexerMockRecorder) AddImageCVEEdges(imagecveedges interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImageCVEEdges", reflect.TypeOf((*MockIndexer)(nil).AddImageCVEEdges), imagecveedges)
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

// DeleteImageCVEEdge mocks base method.
func (m *MockIndexer) DeleteImageCVEEdge(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImageCVEEdge", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImageCVEEdge indicates an expected call of DeleteImageCVEEdge.
func (mr *MockIndexerMockRecorder) DeleteImageCVEEdge(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImageCVEEdge", reflect.TypeOf((*MockIndexer)(nil).DeleteImageCVEEdge), id)
}

// DeleteImageCVEEdges mocks base method.
func (m *MockIndexer) DeleteImageCVEEdges(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImageCVEEdges", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImageCVEEdges indicates an expected call of DeleteImageCVEEdges.
func (mr *MockIndexerMockRecorder) DeleteImageCVEEdges(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImageCVEEdges", reflect.TypeOf((*MockIndexer)(nil).DeleteImageCVEEdges), ids)
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
