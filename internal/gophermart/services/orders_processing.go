package services

import (
	"errors"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
)

type OrdersProcessorInterface interface {
	AddOrder(orderID int, userID int) error
	ValidateOrderID(s int) error
}

type OrdersProcessor struct {
	OrdersStorage storage.OrdersStorage
}

func NewOrdersProcessor(ordersStorage storage.OrdersStorage) *OrdersProcessor {
	return &OrdersProcessor{
		ordersStorage,
	}
}

func (o *OrdersProcessor) AddOrder(orderID int, userID int) error {
	existingOrder, err := o.OrdersStorage.FindByID(orderID)
	if err != nil {
		return err
	}

	if existingOrder.ID != 0 {
		return NewOrderAlreadyAddedError(existingOrder)
	}

	_, err = o.OrdersStorage.CreateNew(orderID, userID)
	if err != nil {
		return err
	}

	// TODO: order processing

	return nil
}

func (o *OrdersProcessor) ValidateOrderID(orderID int) error {
	if orderID <= 0 {
		return errors.New("order ID should be greater than zero")
	}

	if !validLuhn(orderID) {
		return errors.New("order ID is not validated by Luhn algorithm")
	}

	return nil
}
