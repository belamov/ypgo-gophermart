package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrdersRepositoryTestSuite struct {
	suite.Suite
	ordersRepository *OrdersRepository
}

func TestOrdersRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersRepositoryTestSuite))
}

func (s *OrdersRepositoryTestSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db_accrual:5432/accrual?sslmode=disable"
	}
	ordersRepository, err := NewOrdersRepository(dsn)
	require.NoError(s.T(), err)
	s.ordersRepository = ordersRepository

	err = RunMigrations(dsn)
	require.NoError(s.T(), err)

	err = storage.RunMigrations(dsn)
	require.NoError(s.T(), err)
}

func (s *OrdersRepositoryTestSuite) TearDownTest() {
	conn, err := s.ordersRepository.pool.Acquire(context.Background())
	require.NoError(s.T(), err)
	_, _ = conn.Exec(context.Background(), "truncate table accrual_orders cascade")
	conn.Release()
}

func (s *OrdersRepositoryTestSuite) TestCreateNew() {
	orderID := 123
	orderItems := []models.OrderItem{
		{Description: "item 1", Price: 10},
		{Description: "item 2", Price: 20},
	}
	order, err := s.ordersRepository.CreateNew(orderID, orderItems)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), orderID, order.ID)
	assert.Equal(s.T(), orderItems, order.Items)

	order, err = s.ordersRepository.CreateNew(orderID, orderItems)
	var nue *NotUniqueError
	assert.ErrorAs(s.T(), err, &nue)
	assert.Equal(s.T(), models.Order{}, order)
}

func (s *OrdersRepositoryTestSuite) TestAddAccrual() {
	orderID := 123
	orderItems := []models.OrderItem{
		{Description: "item 1", Price: 10},
		{Description: "item 2", Price: 20},
	}
	order, err := s.ordersRepository.CreateNew(orderID, orderItems)
	assert.NoError(s.T(), err)

	err = s.ordersRepository.AddAccrual(order.ID, 10.5)
	assert.NoError(s.T(), err)

	updatedOrder, err := s.ordersRepository.GetOrder(orderID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), models.OrderStatusProcessed, updatedOrder.Status)
	assert.Equal(s.T(), 10.5, updatedOrder.Accrual)
}

func (s *OrdersRepositoryTestSuite) TestChangeStatus() {
	orderID := 123
	orderItems := []models.OrderItem{
		{Description: "item 1", Price: 10},
		{Description: "item 2", Price: 20},
	}
	order, err := s.ordersRepository.CreateNew(orderID, orderItems)
	assert.NoError(s.T(), err)

	err = s.ordersRepository.ChangeStatus(order.ID, models.OrderStatusProcessed)
	assert.NoError(s.T(), err)

	updatedOrder, err := s.ordersRepository.GetOrder(orderID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), models.OrderStatusProcessed, updatedOrder.Status)
	assert.Equal(s.T(), 0.0, updatedOrder.Accrual)
}

func (s *OrdersRepositoryTestSuite) TestGetOrder() {
	newOrder := models.Order{
		ID:        1,
		CreatedAt: time.Now(),
		Status:    models.OrderStatusNew,
		Accrual:   14.3,
		Items:     nil,
	}
	s.add(newOrder)

	order, err := s.ordersRepository.GetOrder(newOrder.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), newOrder.ID, order.ID)
	assert.Equal(s.T(), newOrder.Accrual, order.Accrual)
	assert.Equal(s.T(), newOrder.Status, order.Status)
}

func (s *OrdersRepositoryTestSuite) TestGetOrdersForProcessing() {
	newOrder := models.Order{
		ID:     1,
		Status: models.OrderStatusNew,
	}
	processingOrder := models.Order{
		ID:     2,
		Status: models.OrderStatusProcessing,
	}
	processedOrder := models.Order{
		ID:      3,
		Status:  models.OrderStatusProcessed,
		Accrual: 0,
	}
	erroredOrder := models.Order{
		ID:      4,
		Status:  models.OrderStatusError,
		Accrual: 0,
	}
	s.add(newOrder)
	s.add(processingOrder)
	s.add(processedOrder)
	s.add(erroredOrder)

	ordersToProcess, err := s.ordersRepository.GetOrdersForProcessing()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), ordersToProcess, 1)
	assert.Equal(s.T(), erroredOrder.ID, ordersToProcess[0].ID)
}

func (s *OrdersRepositoryTestSuite) add(order models.Order) {
	conn, err := s.ordersRepository.pool.Acquire(context.Background())
	require.NoError(s.T(), err)

	_, err = conn.Exec(
		context.Background(),
		"insert into accrual_orders (id, status, accrual, created_at) values ($1, $2, $3, $4)",
		order.ID,
		order.Status,
		order.Accrual,
		order.CreatedAt,
	)

	conn.Release()

	assert.NoError(s.T(), err)
}
