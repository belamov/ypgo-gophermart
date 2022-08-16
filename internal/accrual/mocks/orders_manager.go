// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/belamov/ypgo-gophermart/internal/accrual/services (interfaces: OrderManagementInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/belamov/ypgo-gophermart/internal/accrual/models"
	gomock "github.com/golang/mock/gomock"
)

// MockOrderManagementInterface is a mock of OrderManagementInterface interface.
type MockOrderManagementInterface struct {
	ctrl     *gomock.Controller
	recorder *MockOrderManagementInterfaceMockRecorder
}

// MockOrderManagementInterfaceMockRecorder is the mock recorder for MockOrderManagementInterface.
type MockOrderManagementInterfaceMockRecorder struct {
	mock *MockOrderManagementInterface
}

// NewMockOrderManagementInterface creates a new mock instance.
func NewMockOrderManagementInterface(ctrl *gomock.Controller) *MockOrderManagementInterface {
	mock := &MockOrderManagementInterface{ctrl: ctrl}
	mock.recorder = &MockOrderManagementInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderManagementInterface) EXPECT() *MockOrderManagementInterfaceMockRecorder {
	return m.recorder
}

// GetOrderInfo mocks base method.
func (m *MockOrderManagementInterface) GetOrderInfo(arg0 int) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderInfo", arg0)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderInfo indicates an expected call of GetOrderInfo.
func (mr *MockOrderManagementInterfaceMockRecorder) GetOrderInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderInfo", reflect.TypeOf((*MockOrderManagementInterface)(nil).GetOrderInfo), arg0)
}

// RegisterNewOrder mocks base method.
func (m *MockOrderManagementInterface) RegisterNewOrder(arg0 int, arg1 []models.OrderItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterNewOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterNewOrder indicates an expected call of RegisterNewOrder.
func (mr *MockOrderManagementInterfaceMockRecorder) RegisterNewOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterNewOrder", reflect.TypeOf((*MockOrderManagementInterface)(nil).RegisterNewOrder), arg0, arg1)
}

// ValidateOrderID mocks base method.
func (m *MockOrderManagementInterface) ValidateOrderID(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateOrderID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateOrderID indicates an expected call of ValidateOrderID.
func (mr *MockOrderManagementInterfaceMockRecorder) ValidateOrderID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateOrderID", reflect.TypeOf((*MockOrderManagementInterface)(nil).ValidateOrderID), arg0)
}
