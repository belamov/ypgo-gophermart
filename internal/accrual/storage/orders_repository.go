package storage

import (
	"context"
	"log"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/jackc/pgx/v4"
)

type OrdersRepository struct {
	conn *pgx.Conn
}

func (repo *OrdersRepository) Exists(orderID int) (bool, error) {
	var exists bool

	err := repo.conn.QueryRow(
		context.Background(),
		"select exists(select 1 from orders where id = $1)",
		orderID,
	).Scan(&exists)

	return exists, err
}

func NewOrdersRepository(dsn string) (*OrdersRepository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &OrdersRepository{
		conn: conn,
	}, nil
}

func (repo *OrdersRepository) CreateNew(orderID int, items []models.OrderItem) error {
	tx, err := repo.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Default().Println(err)
		}
	}(tx, context.Background())

	err = repo.saveOrder(orderID)
	if err != nil {
		return err
	}

	err = repo.saveOrderItems(orderID, items)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (repo *OrdersRepository) saveOrderItems(orderID int, items []models.OrderItem) error {
	_, err := repo.conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"order_items"},
		[]string{"order_id", "description", "price"},
		pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
			return []interface{}{orderID, items[i].Description, items[i].Price}, nil
		}),
	)
	return err
}

func (repo *OrdersRepository) saveOrder(orderID int) error {
	_, err := repo.conn.Exec(
		context.Background(),
		"insert into orders (id, created_at, status) values ($1, $2, $3)",
		orderID,
		time.Now(),
		models.OrderStatusNew,
	)
	return err
}
