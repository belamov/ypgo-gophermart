package storage

import (
	"context"
	"errors"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type OrdersRepository struct {
	pool *pgxpool.Pool
}

func NewOrdersRepository(dsn string) (*OrdersRepository, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &OrdersRepository{
		pool: pool,
	}, nil
}

func (repo *OrdersRepository) Exists(orderID int) (bool, error) {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return true, err
	}

	var exists bool

	err = conn.QueryRow(
		context.Background(),
		"select exists(select 1 from orders where id = $1)",
		orderID,
	).Scan(&exists)

	conn.Release()

	return exists, err
}

func (repo *OrdersRepository) CreateNew(orderID int, items []models.OrderItem) error {
	conn, err := repo.pool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return err
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Error().Err(err).Msg("unexpected error in rollback transaction")
		}
	}(tx, context.Background())

	err = repo.saveOrder(conn, orderID)
	if err != nil {
		return err
	}

	err = repo.saveOrderItems(conn, orderID, items)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (repo *OrdersRepository) saveOrderItems(conn *pgxpool.Conn, orderID int, items []models.OrderItem) error {
	_, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"order_items"},
		[]string{"order_id", "description", "price"},
		pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
			return []interface{}{orderID, items[i].Description, items[i].Price}, nil
		}),
	)

	return err
}

func (repo *OrdersRepository) saveOrder(conn *pgxpool.Conn, orderID int) error {
	_, err := conn.Exec(
		context.Background(),
		"insert into orders (id, created_at, status) values ($1, $2, $3)",
		orderID,
		time.Now(),
		models.OrderStatusNew,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return NewNotUniqueError("id", err)
		}
	}

	return err
}
