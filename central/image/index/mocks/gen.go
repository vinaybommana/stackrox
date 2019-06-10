// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/image/index (interfaces: Indexer)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	reflect "reflect"
)

// MockIndexer is a mock of Indexer interface
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// AddImage mocks base method
func (m *MockIndexer) AddImage(arg0 *storage.Image) error {
	ret := m.ctrl.Call(m, "AddImage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImage indicates an expected call of AddImage
func (mr *MockIndexerMockRecorder) AddImage(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImage", reflect.TypeOf((*MockIndexer)(nil).AddImage), arg0)
}

// AddImages mocks base method
func (m *MockIndexer) AddImages(arg0 []*storage.Image) error {
	ret := m.ctrl.Call(m, "AddImages", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImages indicates an expected call of AddImages
func (mr *MockIndexerMockRecorder) AddImages(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImages", reflect.TypeOf((*MockIndexer)(nil).AddImages), arg0)
}

// DeleteImage mocks base method
func (m *MockIndexer) DeleteImage(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteImage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImage indicates an expected call of DeleteImage
func (mr *MockIndexerMockRecorder) DeleteImage(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockIndexer)(nil).DeleteImage), arg0)
}

// DeleteImages mocks base method
func (m *MockIndexer) DeleteImages(arg0 []string) error {
	ret := m.ctrl.Call(m, "DeleteImages", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImages indicates an expected call of DeleteImages
func (mr *MockIndexerMockRecorder) DeleteImages(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImages", reflect.TypeOf((*MockIndexer)(nil).DeleteImages), arg0)
}

// GetTxnCount mocks base method
func (m *MockIndexer) GetTxnCount() uint64 {
	ret := m.ctrl.Call(m, "GetTxnCount")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetTxnCount indicates an expected call of GetTxnCount
func (mr *MockIndexerMockRecorder) GetTxnCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxnCount", reflect.TypeOf((*MockIndexer)(nil).GetTxnCount))
}

// ResetIndex mocks base method
func (m *MockIndexer) ResetIndex() error {
	ret := m.ctrl.Call(m, "ResetIndex")
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetIndex indicates an expected call of ResetIndex
func (mr *MockIndexerMockRecorder) ResetIndex() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetIndex", reflect.TypeOf((*MockIndexer)(nil).ResetIndex))
}

// Search mocks base method
func (m *MockIndexer) Search(arg0 *v1.Query) ([]search.Result, error) {
	ret := m.ctrl.Call(m, "Search", arg0)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockIndexerMockRecorder) Search(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIndexer)(nil).Search), arg0)
}

// SetTxnCount mocks base method
func (m *MockIndexer) SetTxnCount(arg0 uint64) error {
	ret := m.ctrl.Call(m, "SetTxnCount", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTxnCount indicates an expected call of SetTxnCount
func (mr *MockIndexerMockRecorder) SetTxnCount(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTxnCount", reflect.TypeOf((*MockIndexer)(nil).SetTxnCount), arg0)
}
