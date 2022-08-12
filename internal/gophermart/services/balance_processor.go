package services

import (
	"sync"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
)

type BalanceProcessorInterface interface {
	RegisterWithdraw(orderID int, userID int, withdrawAmount float64) error
	GetUserTotalAccrualAmount(userID int) (float64, error)
	GetUserTotalWithdrawAmount(userID int) (float64, error)
	GetUserWithdrawals(userID int) ([]models.Withdrawal, error)
	AddAccrual(order models.Order, accrual float64) error
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

	currentBalance, err := b.getUserBalance(userID)
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

func (b *BalanceProcessor) GetUserWithdrawals(userID int) ([]models.Withdrawal, error) {
	withdrawals, err := b.balanceStorage.GetUserWithdrawals(userID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (b *BalanceProcessor) getUserBalance(userID int) (float64, error) {
	totalAccrual, err := b.GetUserTotalAccrualAmount(userID)
	if err != nil {
		return 0, err
	}

	totalWithdraws, err := b.GetUserTotalWithdrawAmount(userID)
	if err != nil {
		return 0, err
	}
	return totalAccrual - totalWithdraws, nil
}

func (b *BalanceProcessor) GetUserTotalAccrualAmount(userID int) (float64, error) {
	totalAccrual, err := b.balanceStorage.GetTotalAccrualAmount(userID)
	if err != nil {
		return 0, err
	}
	return totalAccrual, nil
}

func (b *BalanceProcessor) GetUserTotalWithdrawAmount(userID int) (float64, error) {
	totalWithdraws, err := b.balanceStorage.GetTotalWithdrawAmount(userID)
	if err != nil {
		return 0, err
	}
	return totalWithdraws, nil
}

func (b *BalanceProcessor) AddAccrual(order models.Order, accrual float64) error {
	err := b.balanceStorage.AddAccrual(order.ID, accrual)
	if err != nil {
		return err
	}
	return nil
}

func NewBalanceProcessor(balanceStorage storage.BalanceStorage) *BalanceProcessor {
	return &BalanceProcessor{
		balanceStorage: balanceStorage,
	}
}
