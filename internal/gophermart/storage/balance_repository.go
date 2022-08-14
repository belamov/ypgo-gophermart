package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/jackc/pgtype"
)

type BalanceRepository struct {
	pool *pgxpool.Pool
}

func NewBalanceRepository(dsn string) (*BalanceRepository, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &BalanceRepository{
		pool: pool,
	}, nil
}

func (repo *BalanceRepository) AddAccrual(orderID int, accrualAmount float64) error {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return err
	}

	_, err = conn.Exec(
		context.Background(),
		"update orders set status=$1, accrual=$2 where id=$3",
		models.OrderStatusProcessed,
		accrualAmount,
		orderID,
	)

	conn.Release()

	return err
}

func (repo *BalanceRepository) GetUserWithdrawals(userID int) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return nil, err
	}

	rows, err := conn.Query(
		context.Background(),
		"select order_id, user_id, amount, created_at from withdraws where user_id=$1 order by created_at",
		userID,
	)

	conn.Release()

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

func (repo *BalanceRepository) GetTotalAccrualAmount(userID int) (float64, error) {
	result := 0.0

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return result, err
	}

	err = conn.QueryRow(
		context.Background(),
		"select sum(accrual) from orders where created_by=$1 group by created_by",
		userID,
	).Scan(&result)

	conn.Release()

	if err == pgx.ErrNoRows {
		return 0.0, nil
	}

	return result, err
}

func (repo *BalanceRepository) GetTotalWithdrawAmount(userID int) (float64, error) {
	result := 0.0

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return result, err
	}

	err = conn.QueryRow(
		context.Background(),
		"select sum(amount) from withdraws where user_id=$1 group by user_id",
		userID,
	).Scan(&result)

	conn.Release()

	if err == pgx.ErrNoRows {
		return 0.0, nil
	}
	
	return result, err
}

func (repo *BalanceRepository) AddWithdraw(orderID int, userID int, amount float64) error {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return err
	}

	_, err = conn.Exec(
		context.Background(),
		"insert into withdraws (order_id, user_id, amount) values ($1, $2, $3)",
		orderID,
		userID,
		amount,
	)

	conn.Release()

	return err
}
