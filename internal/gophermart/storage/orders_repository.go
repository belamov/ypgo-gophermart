package storage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/jackc/pgx/v4"
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

func (repo *OrdersRepository) ChangeStatus(order models.Order, status models.OrderStatus) error {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return err
	}

	_, err = conn.Exec(
		context.Background(),
		"update orders set status=$1 where id=$2",
		status,
		order.ID,
	)

	conn.Release()

	return err
}

func (repo *OrdersRepository) GetUsersOrders(userID int) ([]models.Order, error) {
	var orders []models.Order

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return nil, err
	}

	rows, err := conn.Query(
		context.Background(),
		"select id, created_by, uploaded_at, status, accrual from orders where created_by=$1 order by uploaded_at desc",
		userID,
	)

	conn.Release()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		model := models.Order{}
		var accrual sql.NullFloat64
		if err = rows.Scan(&model.ID, &model.CreatedBy, &model.UploadedAt, &model.Status, &accrual); err != nil {
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

func (repo *OrdersRepository) FindByID(orderID int) (models.Order, error) {
	var order models.Order
	var accrual sql.NullFloat64

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return order, err
	}

	err = conn.QueryRow(
		context.Background(),
		"select id, created_by, uploaded_at, status, accrual from orders where id=$1",
		orderID,
	).Scan(&order.ID, &order.CreatedBy, &order.UploadedAt, &order.Status, &accrual)

	conn.Release()

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

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Println("couldnt acquire connection from pool:")
		log.Println(err.Error())
		return models.Order{}, err
	}

	_, err = conn.Exec(
		context.Background(),
		"insert into orders (id, created_by, uploaded_at, status) values ($1, $2, $3, $4)",
		order.ID,
		order.CreatedBy,
		order.UploadedAt,
		order.Status,
	)

	conn.Release()

	// todo: handle not unique order id

	return order, err
}
