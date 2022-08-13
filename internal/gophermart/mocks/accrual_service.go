// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/belamov/ypgo-gophermart/internal/gophermart/services (interfaces: AccrualInfoProvider)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAccrualInfoProvider is a mock of AccrualInfoProvider interface.
type MockAccrualInfoProvider struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualInfoProviderMockRecorder
}

// MockAccrualInfoProviderMockRecorder is the mock recorder for MockAccrualInfoProvider.
type MockAccrualInfoProviderMockRecorder struct {
	mock *MockAccrualInfoProvider
}

// NewMockAccrualInfoProvider creates a new mock instance.
func NewMockAccrualInfoProvider(ctrl *gomock.Controller) *MockAccrualInfoProvider {
	mock := &MockAccrualInfoProvider{ctrl: ctrl}
	mock.recorder = &MockAccrualInfoProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualInfoProvider) EXPECT() *MockAccrualInfoProviderMockRecorder {
	return m.recorder
}

// GetAccrualForOrder mocks base method.
func (m *MockAccrualInfoProvider) GetAccrualForOrder(arg0 context.Context, arg1 int) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccrualForOrder", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccrualForOrder indicates an expected call of GetAccrualForOrder.
func (mr *MockAccrualInfoProviderMockRecorder) GetAccrualForOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccrualForOrder", reflect.TypeOf((*MockAccrualInfoProvider)(nil).GetAccrualForOrder), arg0, arg1)
}
