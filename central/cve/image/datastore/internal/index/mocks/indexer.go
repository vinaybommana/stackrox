// Code generated by MockGen. DO NOT EDIT.
// Source: indexer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
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

// AddImageCVE mocks base method.
func (m *MockIndexer) AddImageCVE(cve *storage.ImageCVE) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImageCVE", cve)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImageCVE indicates an expected call of AddImageCVE.
func (mr *MockIndexerMockRecorder) AddImageCVE(cve interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImageCVE", reflect.TypeOf((*MockIndexer)(nil).AddImageCVE), cve)
}

// AddImageCVEs mocks base method.
func (m *MockIndexer) AddImageCVEs(cves []*storage.ImageCVE) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImageCVEs", cves)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImageCVEs indicates an expected call of AddImageCVEs.
func (mr *MockIndexerMockRecorder) AddImageCVEs(cves interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImageCVEs", reflect.TypeOf((*MockIndexer)(nil).AddImageCVEs), cves)
}

// Count mocks base method.
func (m *MockIndexer) Count(q *v1.Query, opts ...blevesearch.SearchOption) (int, error) {
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

// DeleteImageCVE mocks base method.
func (m *MockIndexer) DeleteImageCVE(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImageCVE", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImageCVE indicates an expected call of DeleteImageCVE.
func (mr *MockIndexerMockRecorder) DeleteImageCVE(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImageCVE", reflect.TypeOf((*MockIndexer)(nil).DeleteImageCVE), id)
}

// DeleteImageCVEs mocks base method.
func (m *MockIndexer) DeleteImageCVEs(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImageCVEs", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImageCVEs indicates an expected call of DeleteImageCVEs.
func (mr *MockIndexerMockRecorder) DeleteImageCVEs(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImageCVEs", reflect.TypeOf((*MockIndexer)(nil).DeleteImageCVEs), ids)
}

// Search mocks base method.
func (m *MockIndexer) Search(q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
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
