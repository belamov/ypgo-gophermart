package storage

import (
	"context"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type OrdersRepository struct {
	conn *pgx.Conn
}

func (repo *OrdersRepository) GetUsersOrders(userID int) ([]models.Order, error) {
	// TODO implement me
	panic("implement me")
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
	var accrual pgtype.Float8
	err := repo.conn.QueryRow(
		context.Background(),
		"select id, created_by, uploaded_at, status, accrual from orders where id=$1",
		orderID,
	).Scan(&order.ID, &order.CreatedBy, &order.UploadedAt, &order.Status, &accrual)
	order.Accrual = accrual.Float

	if err == pgx.ErrNoRows {
		return models.Order{}, nil
	}

	return order, err
}

func (repo *OrdersRepository) CreateNew(orderID int, userID int) (models.Order, error) {
	order := models.Order{
		ID:         orderID,
		CreatedBy:  userID,
		UploadedAt: time.Now(),
		Status:     models.OrderStatusNew,
	}

	_, err := repo.conn.Exec(
		context.Background(),
		"insert into orders (id, created_by, uploaded_at, status) values ($1, $2, $3, $4)",
		order.ID,
		order.CreatedBy,
		order.UploadedAt,
		order.Status,
	)

	return order, err
}
