// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/belamov/ypgo-gophermart/internal/gophermart/services (interfaces: OrdersProcessorInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	gomock "github.com/golang/mock/gomock"
)

// MockOrdersProcessorInterface is a mock of OrdersProcessorInterface interface.
type MockOrdersProcessorInterface struct {
	ctrl     *gomock.Controller
	recorder *MockOrdersProcessorInterfaceMockRecorder
}

// MockOrdersProcessorInterfaceMockRecorder is the mock recorder for MockOrdersProcessorInterface.
type MockOrdersProcessorInterfaceMockRecorder struct {
	mock *MockOrdersProcessorInterface
}

// NewMockOrdersProcessorInterface creates a new mock instance.
func NewMockOrdersProcessorInterface(ctrl *gomock.Controller) *MockOrdersProcessorInterface {
	mock := &MockOrdersProcessorInterface{ctrl: ctrl}
	mock.recorder = &MockOrdersProcessorInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrdersProcessorInterface) EXPECT() *MockOrdersProcessorInterfaceMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockOrdersProcessorInterface) AddOrder(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockOrdersProcessorInterfaceMockRecorder) AddOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockOrdersProcessorInterface)(nil).AddOrder), arg0, arg1)
}

// GetUsersOrders mocks base method.
func (m *MockOrdersProcessorInterface) GetUsersOrders(arg0 int) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersOrders", arg0)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersOrders indicates an expected call of GetUsersOrders.
func (mr *MockOrdersProcessorInterfaceMockRecorder) GetUsersOrders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersOrders", reflect.TypeOf((*MockOrdersProcessorInterface)(nil).GetUsersOrders), arg0)
}

// ValidateOrderID mocks base method.
func (m *MockOrdersProcessorInterface) ValidateOrderID(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateOrderID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateOrderID indicates an expected call of ValidateOrderID.
func (mr *MockOrdersProcessorInterfaceMockRecorder) ValidateOrderID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateOrderID", reflect.TypeOf((*MockOrdersProcessorInterface)(nil).ValidateOrderID), arg0)
}
