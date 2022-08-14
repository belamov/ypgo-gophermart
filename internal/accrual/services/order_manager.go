package services

import (
	"errors"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"

	"github.com/belamov/ypgo-gophermart/internal"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
)

type OrderManagementInterface interface {
	ValidateOrderID(orderID int) error
	RegisterNewOrder(orderID int, orderItems []models.OrderItem) error
}

var ErrOrderIsAlreadyRegistered = errors.New("order is already registered")

type OrderManager struct {
	orderStorage storage.OrdersStorage
}

func NewOrderManager(orderStorage storage.OrdersStorage) *OrderManager {
	return &OrderManager{orderStorage: orderStorage}
}

func (o *OrderManager) RegisterNewOrder(orderID int, orderItems []models.OrderItem) error {
	isOrderRegistered, err := o.orderStorage.IsRegistered(orderID)
	if err != nil {
		return err
	}

	if isOrderRegistered {
		return ErrOrderIsAlreadyRegistered
	}

	err = o.orderStorage.RegisterOrder(orderID, orderItems)
	if err != nil {
		return err
	}

	// todo: order processing

	return nil
}

func (o *OrderManager) ValidateOrderID(orderID int) error {
	if orderID <= 0 {
		return errors.New("order ID should be greater than zero")
	}

	if !internal.ValidLuhn(orderID) {
		return errors.New("order ID is not validated by Luhn algorithm")
	}

	return nil
}
