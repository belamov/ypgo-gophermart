package storage

import (
	"context"
	"database/sql"
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

func (repo *OrdersRepository) CreateNew(orderID int, items []models.OrderItem) (models.Order, error) {
	conn, err := repo.pool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return models.Order{}, err
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return models.Order{}, err
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
		return models.Order{}, err
	}

	err = repo.saveOrderItems(conn, orderID, items)
	if err != nil {
		return models.Order{}, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return models.Order{}, err
	}

	return models.Order{ID: orderID, Items: items}, nil
}

func (repo *OrdersRepository) AddAccrual(orderID int, accrual float64) error {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return err
	}

	_, err = conn.Exec(
		context.Background(),
		"update orders set accrual=$1 where id=$2",
		accrual,
		orderID,
	)

	conn.Release()

	return err
}

func (repo *OrdersRepository) ChangeStatus(orderID int, status models.OrderStatus) error {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return err
	}

	_, err = conn.Exec(
		context.Background(),
		"update orders set status=$1 where id=$2",
		status,
		orderID,
	)

	conn.Release()

	return err
}

func (repo *OrdersRepository) GetOrdersForProcessing() ([]models.Order, error) {
	var orders []models.Order

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return nil, err
	}

	rows, err := conn.Query(
		context.Background(),
		"select id, created_at, status, accrual from orders where status=$1 order by created_at",
		models.OrderStatusError,
	)

	conn.Release()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		model := models.Order{}
		var accrual sql.NullFloat64
		if err = rows.Scan(&model.ID, &model.CreatedAt, &model.Status, &accrual); err != nil {
			return nil, err
		}
		model.Accrual = accrual.Float64
		orders = append(orders, model)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return orders, nil
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
