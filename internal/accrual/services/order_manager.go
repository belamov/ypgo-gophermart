package services

import (
	"errors"

	"github.com/belamov/ypgo-gophermart/internal"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
)

type OrderManagementInterface interface {
	ValidateOrderID(orderID int) error
	RegisterNewOrder(orderID int, orderItems []models.OrderItem) error
}

var ErrOrderIsAlreadyRegistered = errors.New("order is already registered")

type OrderManager struct{}

func (o OrderManager) ValidateOrderID(orderID int) error {
	if orderID <= 0 {
		return errors.New("order ID should be greater than zero")
	}

	if !internal.ValidLuhn(orderID) {
		return errors.New("order ID is not validated by Luhn algorithm")
	}

	return nil
}
