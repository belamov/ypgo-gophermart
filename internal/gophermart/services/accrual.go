package services

import "errors"

type AccrualInfoProvider interface {
	GetAccrualForOrder(orderID int) (float64, error)
}

var (
	ErrOrderIsNotYetProceeded = errors.New(`order is not yet proceeded`)
	ErrInvalidOrderForAccrual = errors.New(`this order is invalid`)
)
