// Code generated by MockGen. DO NOT EDIT.
// Source: searcher.go

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

// MockSearcher is a mock of Searcher interface.
type MockSearcher struct {
	ctrl     *gomock.Controller
	recorder *MockSearcherMockRecorder
}

// MockSearcherMockRecorder is the mock recorder for MockSearcher.
type MockSearcherMockRecorder struct {
	mock *MockSearcher
}

// NewMockSearcher creates a new mock instance.
func NewMockSearcher(ctrl *gomock.Controller) *MockSearcher {
	mock := &MockSearcher{ctrl: ctrl}
	mock.recorder = &MockSearcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearcher) EXPECT() *MockSearcherMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockSearcher) Count(ctx context.Context, q *auxpb.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, q)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockSearcherMockRecorder) Count(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockSearcher)(nil).Count), ctx, q)
}

// Search mocks base method.
func (m *MockSearcher) Search(ctx context.Context, q *auxpb.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockSearcherMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockSearcher)(nil).Search), ctx, q)
}

// SearchDeployments mocks base method.
func (m *MockSearcher) SearchDeployments(ctx context.Context, q *auxpb.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchDeployments", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchDeployments indicates an expected call of SearchDeployments.
func (mr *MockSearcherMockRecorder) SearchDeployments(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchDeployments", reflect.TypeOf((*MockSearcher)(nil).SearchDeployments), ctx, q)
}

// SearchListDeployments mocks base method.
func (m *MockSearcher) SearchListDeployments(ctx context.Context, q *auxpb.Query) ([]*storage.ListDeployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchListDeployments", ctx, q)
	ret0, _ := ret[0].([]*storage.ListDeployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchListDeployments indicates an expected call of SearchListDeployments.
func (mr *MockSearcherMockRecorder) SearchListDeployments(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchListDeployments", reflect.TypeOf((*MockSearcher)(nil).SearchListDeployments), ctx, q)
}

// SearchRawDeployments mocks base method.
func (m *MockSearcher) SearchRawDeployments(ctx context.Context, q *auxpb.Query) ([]*storage.Deployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawDeployments", ctx, q)
	ret0, _ := ret[0].([]*storage.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawDeployments indicates an expected call of SearchRawDeployments.
func (mr *MockSearcherMockRecorder) SearchRawDeployments(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawDeployments", reflect.TypeOf((*MockSearcher)(nil).SearchRawDeployments), ctx, q)
}
