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

	if err = runMigrations(dsn, "file://./migrations"); err != nil {
		return nil, err
	}

	return &OrdersRepository{
		conn: conn,
	}, nil
}

func (o *OrdersRepository) FindByID(orderID int) (models.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (o *OrdersRepository) CreateNew(orderID int, userID int) (models.Order, error) {
	// TODO implement me
	panic("implement me")
}
