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

	return &BalanceRepository{
		conn: conn,
	}, nil
}

func (repo *BalanceRepository) GetTotalAccrual(userID int) (float64, error) {
	result := 0.0
	err := repo.conn.QueryRow(
		context.Background(),
		"select sum(accrual) from orders where created_by=$1 group by created_by",
		userID,
	).Scan(&result)

	return result, err
}

func (repo *BalanceRepository) GetTotalWithdraws(userID int) (float64, error) {
	result := 0.0
	err := repo.conn.QueryRow(
		context.Background(),
		"select sum(amount) from withdraws where user_id=$1 group by user_id",
		userID,
	).Scan(&result)

	return result, err
}

func (repo *BalanceRepository) AddWithdraw(orderID int, userID int, amount float64) error {
	_, err := repo.conn.Exec(
		context.Background(),
		"insert into withdraws (order_id, user_id, amount) values ($1, $2, $3)",
		orderID,
		userID,
		amount,
	)

	return err
}
