// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/belamov/ypgo-gophermart/internal/gophermart/storage (interfaces: OrdersStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	gomock "github.com/golang/mock/gomock"
)

// MockOrdersStorage is a mock of OrdersStorage interface.
type MockOrdersStorage struct {
	ctrl     *gomock.Controller
	recorder *MockOrdersStorageMockRecorder
}

// MockOrdersStorageMockRecorder is the mock recorder for MockOrdersStorage.
type MockOrdersStorageMockRecorder struct {
	mock *MockOrdersStorage
}

// NewMockOrdersStorage creates a new mock instance.
func NewMockOrdersStorage(ctrl *gomock.Controller) *MockOrdersStorage {
	mock := &MockOrdersStorage{ctrl: ctrl}
	mock.recorder = &MockOrdersStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrdersStorage) EXPECT() *MockOrdersStorageMockRecorder {
	return m.recorder
}

// CreateNew mocks base method.
func (m *MockOrdersStorage) CreateNew(arg0, arg1 int) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNew", arg0, arg1)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNew indicates an expected call of CreateNew.
func (mr *MockOrdersStorageMockRecorder) CreateNew(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNew", reflect.TypeOf((*MockOrdersStorage)(nil).CreateNew), arg0, arg1)
}

// FindByID mocks base method.
func (m *MockOrdersStorage) FindByID(arg0 int) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockOrdersStorageMockRecorder) FindByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockOrdersStorage)(nil).FindByID), arg0)
}