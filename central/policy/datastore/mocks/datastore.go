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

// AddPolicy mocks base method.
func (m *MockDataStore) AddPolicy(arg0 context.Context, arg1 *storage.Policy) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPolicy", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPolicy indicates an expected call of AddPolicy.
func (mr *MockDataStoreMockRecorder) AddPolicy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPolicy", reflect.TypeOf((*MockDataStore)(nil).AddPolicy), arg0, arg1)
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

// GetAllPolicies mocks base method.
func (m *MockDataStore) GetAllPolicies(ctx context.Context) ([]*storage.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPolicies", ctx)
	ret0, _ := ret[0].([]*storage.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPolicies indicates an expected call of GetAllPolicies.
func (mr *MockDataStoreMockRecorder) GetAllPolicies(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPolicies", reflect.TypeOf((*MockDataStore)(nil).GetAllPolicies), ctx)
}

// GetPolicies mocks base method.
func (m *MockDataStore) GetPolicies(ctx context.Context, ids []string) ([]*storage.Policy, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicies", ctx, ids)
	ret0, _ := ret[0].([]*storage.Policy)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPolicies indicates an expected call of GetPolicies.
func (mr *MockDataStoreMockRecorder) GetPolicies(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicies", reflect.TypeOf((*MockDataStore)(nil).GetPolicies), ctx, ids)
}

// GetPolicy mocks base method.
func (m *MockDataStore) GetPolicy(ctx context.Context, id string) (*storage.Policy, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicy", ctx, id)
	ret0, _ := ret[0].(*storage.Policy)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPolicy indicates an expected call of GetPolicy.
func (mr *MockDataStoreMockRecorder) GetPolicy(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicy", reflect.TypeOf((*MockDataStore)(nil).GetPolicy), ctx, id)
}

// GetPolicyByName mocks base method.
func (m *MockDataStore) GetPolicyByName(ctx context.Context, name string) (*storage.Policy, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicyByName", ctx, name)
	ret0, _ := ret[0].(*storage.Policy)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPolicyByName indicates an expected call of GetPolicyByName.
func (mr *MockDataStoreMockRecorder) GetPolicyByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicyByName", reflect.TypeOf((*MockDataStore)(nil).GetPolicyByName), ctx, name)
}

// ImportPolicies mocks base method.
func (m *MockDataStore) ImportPolicies(ctx context.Context, policies []*storage.Policy, overwrite bool) ([]*v1.ImportPolicyResponse, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportPolicies", ctx, policies, overwrite)
	ret0, _ := ret[0].([]*v1.ImportPolicyResponse)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ImportPolicies indicates an expected call of ImportPolicies.
func (mr *MockDataStoreMockRecorder) ImportPolicies(ctx, policies, overwrite interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportPolicies", reflect.TypeOf((*MockDataStore)(nil).ImportPolicies), ctx, policies, overwrite)
}

// RemovePolicy mocks base method.
func (m *MockDataStore) RemovePolicy(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePolicy", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePolicy indicates an expected call of RemovePolicy.
func (mr *MockDataStoreMockRecorder) RemovePolicy(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePolicy", reflect.TypeOf((*MockDataStore)(nil).RemovePolicy), ctx, id)
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

// SearchPolicies mocks base method.
func (m *MockDataStore) SearchPolicies(ctx context.Context, q *auxpb.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchPolicies", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchPolicies indicates an expected call of SearchPolicies.
func (mr *MockDataStoreMockRecorder) SearchPolicies(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchPolicies", reflect.TypeOf((*MockDataStore)(nil).SearchPolicies), ctx, q)
}

// SearchRawPolicies mocks base method.
func (m *MockDataStore) SearchRawPolicies(ctx context.Context, q *auxpb.Query) ([]*storage.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawPolicies", ctx, q)
	ret0, _ := ret[0].([]*storage.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawPolicies indicates an expected call of SearchRawPolicies.
func (mr *MockDataStoreMockRecorder) SearchRawPolicies(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawPolicies", reflect.TypeOf((*MockDataStore)(nil).SearchRawPolicies), ctx, q)
}

// UpdatePolicy mocks base method.
func (m *MockDataStore) UpdatePolicy(arg0 context.Context, arg1 *storage.Policy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePolicy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePolicy indicates an expected call of UpdatePolicy.
func (mr *MockDataStoreMockRecorder) UpdatePolicy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePolicy", reflect.TypeOf((*MockDataStore)(nil).UpdatePolicy), arg0, arg1)
}
