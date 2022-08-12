package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/jackc/pgtype"
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

	return &OrdersRepository{
		conn: conn,
	}, nil
}

func (repo *OrdersRepository) ChangeStatus(order models.Order, status models.OrderStatus) error {
	// TODO implement me
	panic("implement me")
}

func (repo *OrdersRepository) GetUsersOrders(userID int) ([]models.Order, error) {
	var orders []models.Order

	rows, err := repo.conn.Query(
		context.Background(),
		"select id, created_by, uploaded_at, status, accrual from orders where created_by=$1 order by uploaded_at desc",
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		model := models.Order{}
		var accrual pgtype.Float8
		if err = rows.Scan(&model.ID, &model.CreatedBy, &model.UploadedAt, &model.Status, &accrual); err != nil {
			return nil, err
		}
		model.Accrual = accrual.Float
		orders = append(orders, model)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return orders, nil
}

func (repo *OrdersRepository) FindByID(orderID int) (models.Order, error) {
	var order models.Order
	var accrual sql.NullFloat64
	err := repo.conn.QueryRow(
		context.Background(),
		"select id, created_by, uploaded_at, status, accrual from orders where id=$1",
		orderID,
	).Scan(&order.ID, &order.CreatedBy, &order.UploadedAt, &order.Status, &accrual)
	accrualFloat := accrual.Float64
	order.Accrual = accrualFloat

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

	// todo: handle not unique order id

	return order, err
}
