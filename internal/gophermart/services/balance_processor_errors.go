package services

import (
	"fmt"
)

type InsufficientBalanceError struct {
	Err            error
	Balance        float64
	WithdrawAmount float64
}

func (err *InsufficientBalanceError) Error() string {
	return fmt.Sprintf("can't register withdraw: required %v, only got %v", err.WithdrawAmount, err.Balance)
}

func (err *InsufficientBalanceError) Unwrap() error {
	return err.Err
}

func NewInsufficientBalanceError(balance float64, withdrawAmount float64) error {
	return &InsufficientBalanceError{
		Balance:        balance,
		WithdrawAmount: withdrawAmount,
	}
}
