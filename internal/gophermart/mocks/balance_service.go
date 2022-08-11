// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/belamov/ypgo-gophermart/internal/gophermart/services (interfaces: BalanceProcessorInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBalanceProcessorInterface is a mock of BalanceProcessorInterface interface.
type MockBalanceProcessorInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceProcessorInterfaceMockRecorder
}

// MockBalanceProcessorInterfaceMockRecorder is the mock recorder for MockBalanceProcessorInterface.
type MockBalanceProcessorInterfaceMockRecorder struct {
	mock *MockBalanceProcessorInterface
}

// NewMockBalanceProcessorInterface creates a new mock instance.
func NewMockBalanceProcessorInterface(ctrl *gomock.Controller) *MockBalanceProcessorInterface {
	mock := &MockBalanceProcessorInterface{ctrl: ctrl}
	mock.recorder = &MockBalanceProcessorInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalanceProcessorInterface) EXPECT() *MockBalanceProcessorInterfaceMockRecorder {
	return m.recorder
}

// RegisterWithdraw mocks base method.
func (m *MockBalanceProcessorInterface) RegisterWithdraw(arg0, arg1 int, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterWithdraw", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterWithdraw indicates an expected call of RegisterWithdraw.
func (mr *MockBalanceProcessorInterfaceMockRecorder) RegisterWithdraw(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterWithdraw", reflect.TypeOf((*MockBalanceProcessorInterface)(nil).RegisterWithdraw), arg0, arg1, arg2)
}