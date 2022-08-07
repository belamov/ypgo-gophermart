package services

import "errors"

type OrdersProcessorInterface interface {
	AddOrder(orderID int, userID int) error
	ValidateOrderID(s int) error
}

type OrdersProcessor struct{}

func (o *OrdersProcessor) AddOrder(orderID int, userID int) error {
	// TODO implement me
	panic("implement me")
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

func NewOrdersProcessor() *OrdersProcessor {
	return &OrdersProcessor{}
}
