// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/license/manager (interfaces: LicenseEventListener)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	reflect "reflect"
)

// MockLicenseEventListener is a mock of LicenseEventListener interface
type MockLicenseEventListener struct {
	ctrl     *gomock.Controller
	recorder *MockLicenseEventListenerMockRecorder
}

// MockLicenseEventListenerMockRecorder is the mock recorder for MockLicenseEventListener
type MockLicenseEventListenerMockRecorder struct {
	mock *MockLicenseEventListener
}

// NewMockLicenseEventListener creates a new mock instance
func NewMockLicenseEventListener(ctrl *gomock.Controller) *MockLicenseEventListener {
	mock := &MockLicenseEventListener{ctrl: ctrl}
	mock.recorder = &MockLicenseEventListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLicenseEventListener) EXPECT() *MockLicenseEventListenerMockRecorder {
	return m.recorder
}

// OnActiveLicenseChanged mocks base method
func (m *MockLicenseEventListener) OnActiveLicenseChanged(arg0, arg1 *v1.LicenseInfo) {
	m.ctrl.Call(m, "OnActiveLicenseChanged", arg0, arg1)
}

// OnActiveLicenseChanged indicates an expected call of OnActiveLicenseChanged
func (mr *MockLicenseEventListenerMockRecorder) OnActiveLicenseChanged(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnActiveLicenseChanged", reflect.TypeOf((*MockLicenseEventListener)(nil).OnActiveLicenseChanged), arg0, arg1)
}
