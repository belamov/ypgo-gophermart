package storage

import (
	"context"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/jackc/pgx/v4"
)

type OrdersRepository struct {
	conn *pgx.Conn
}

func NewOrdersRepository(dsn string) (*OrdersRepository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = runMigrations(dsn); err != nil {
		return nil, err
	}

	return &OrdersRepository{
		conn: conn,
	}, nil
}

func (repo *OrdersRepository) FindByID(orderID int) (models.Order, error) {
	var order models.Order
	err := repo.conn.QueryRow(
		context.Background(),
		"select id, created_by from orders where id=$1",
		orderID,
	).Scan(&order.ID, &order.CreatedBy)

	if err == pgx.ErrNoRows {
		return models.Order{}, nil
	}

	return order, err
}

func (repo *OrdersRepository) CreateNew(orderID int, userID int) (models.Order, error) {
	order := models.Order{
		ID:        orderID,
		CreatedBy: userID,
	}

	_, err := repo.conn.Exec(
		context.Background(),
		"insert into orders (id, created_by) values ($1, $2)",
		order.ID,
		order.CreatedBy,
	)

	return order, err
}
