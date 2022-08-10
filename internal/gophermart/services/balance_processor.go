package services

import (
	"sync"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
)

type BalanceProcessorInterface interface {
	RegisterWithdraw(orderID int, userID int, withdrawAmount float64) error
}

type BalanceProcessor struct {
	balanceStorage storage.BalanceStorage
	mu             sync.Mutex
}

func (b *BalanceProcessor) RegisterWithdraw(orderID int, userID int, withdrawAmount float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	// todo: lock only for given user id
	// https://pkg.go.dev/golang.org/x/sync@v0.0.0-20220722155255-886fb9371eb4/singleflight

	currentBalance, err := b.getCurrentUserBalance(userID)
	if err != nil {
		return err
	}

	if currentBalance < withdrawAmount {
		return NewInsufficientBalanceError(currentBalance, withdrawAmount)
	}

	err = b.balanceStorage.AddWithdraw(orderID, userID, withdrawAmount)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceProcessor) getCurrentUserBalance(userID int) (float64, error) {
	totalAccrual, err := b.getTotalAccrual(userID)
	if err != nil {
		return 0, err
	}

	totalWithdraws, err := b.getTotalWithdraws(userID)
	if err != nil {
		return 0, err
	}
	return totalAccrual - totalWithdraws, nil
}

func (b *BalanceProcessor) getTotalAccrual(userID int) (float64, error) {
	totalAccrual, err := b.balanceStorage.GetTotalAccrual(userID)
	if err != nil {
		return 0, err
	}
	return totalAccrual, nil
}

func (b *BalanceProcessor) getTotalWithdraws(userID int) (float64, error) {
	totalWithdraws, err := b.balanceStorage.GetTotalWithdraws(userID)
	if err != nil {
		return 0, err
	}
	return totalWithdraws, nil
}

func NewBalanceProcessor(balanceStorage storage.BalanceStorage) *BalanceProcessor {
	return &BalanceProcessor{
		balanceStorage: balanceStorage,
	}
}
