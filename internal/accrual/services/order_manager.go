package services

import (
	"errors"

	"github.com/belamov/ypgo-gophermart/internal"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
)

type OrderManagementInterface interface {
	ValidateOrderID(orderID int) error
	RegisterNewOrder(orderID int, orderItems []models.OrderItem) error
}

var ErrOrderIsAlreadyRegistered = errors.New("order is already registered")

type OrderManager struct {
	orderStorage     storage.OrdersStorage
	accrualProcessor AccrualProcessor
}

func NewOrderManager(orderStorage storage.OrdersStorage) *OrderManager {
	return &OrderManager{orderStorage: orderStorage}
}

func (o *OrderManager) RegisterNewOrder(orderID int, orderItems []models.OrderItem) error {
	order, err := o.orderStorage.CreateNew(orderID, orderItems)
	var errNotUnique *storage.NotUniqueError
	if errors.As(err, &errNotUnique) {
		return ErrOrderIsAlreadyRegistered
	}
	if err != nil {
		return err
	}

	go o.accrualProcessor.RegisterOrderForProcessing(order)

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
