package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type BalanceRepository struct {
	conn *pgx.Conn
}

func NewBalanceRepository(dsn string) (*BalanceRepository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = runMigrations(dsn); err != nil {
		return nil, err
	}

	return &BalanceRepository{
		conn: conn,
	}, nil
}

func (repo *BalanceRepository) GetTotalAccrual(userID int) (float64, error) {
	// TODO implement me
	panic("implement me")
}

func (repo *BalanceRepository) GetTotalWithdraws(userID int) (float64, error) {
	// TODO implement me
	panic("implement me")
}

func (repo *BalanceRepository) AddWithdraw(orderID int, userID int, amount float64) error {
	// TODO implement me
	panic("implement me")
}
