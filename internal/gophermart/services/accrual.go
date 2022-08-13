package services

import (
	"context"
	"errors"
)

type AccrualInfoProvider interface {
	GetAccrualForOrder(ctx context.Context, orderID int) (float64, error)
}

var (
	ErrOrderIsNotYetProceeded = errors.New(`order is not yet proceeded`)
	ErrInvalidOrderForAccrual = errors.New(`this order is invalid`)
)
