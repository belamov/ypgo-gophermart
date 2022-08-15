package services

import (
	"fmt"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
)

type OrderAlreadyAddedError struct {
	Order models.Order
}

func (err *OrderAlreadyAddedError) Error() string {
	return fmt.Sprintf("Order alredy added: %+v", err.Order)
}

func (err *OrderAlreadyAddedError) Unwrap() error {
	return nil
}

func NewOrderAlreadyAddedError(order models.Order) error {
	return &OrderAlreadyAddedError{
		Order: order,
	}
}
