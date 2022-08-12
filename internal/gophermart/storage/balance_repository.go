package storage

import (
	"context"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type BalanceRepository struct {
	conn *pgx.Conn
}

func (repo *BalanceRepository) AddAccrual(orderID int, accrualAmount float64) error {
	_, err := repo.conn.Exec(
		context.Background(),
		"update orders set status=$1, accrual=$2 where id=$3",
		models.OrderStatusProcessed,
		accrualAmount,
		orderID,
	)

	return err
}

func (repo *BalanceRepository) GetUserWithdrawals(userID int) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	rows, err := repo.conn.Query(
		context.Background(),
		"select order_id, user_id, amount, created_at from withdraws where user_id=$1 order by created_at",
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		model := models.Withdrawal{}
		var amount pgtype.Float8
		var createdAt pgtype.Timestamp
		if err = rows.Scan(&model.OrderID, &model.UserID, &amount, &createdAt); err != nil {
			return nil, err
		}
		model.WithdrawalAmount = amount.Float
		model.CreatedAt = createdAt.Time
		withdrawals = append(withdrawals, model)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return withdrawals, nil
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

func (repo *BalanceRepository) GetTotalAccrualAmount(userID int) (float64, error) {
	result := 0.0
	err := repo.conn.QueryRow(
		context.Background(),
		"select sum(accrual) from orders where created_by=$1 group by created_by",
		userID,
	).Scan(&result)

	return result, err
}

func (repo *BalanceRepository) GetTotalWithdrawAmount(userID int) (float64, error) {
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
