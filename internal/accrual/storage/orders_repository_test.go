package storage

import (
	"context"
	"os"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
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
}

func (s *OrdersRepositoryTestSuite) TearDownTest() {
	conn, err := s.ordersRepository.pool.Acquire(context.Background())
	require.NoError(s.T(), err)
	_, _ = conn.Exec(context.Background(), "truncate table orders cascade")
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
		"insert into orders (id, status, accrual) values ($1, $2, $3)",
		order.ID,
		order.Status,
		order.Accrual,
	)

	conn.Release()

	assert.NoError(s.T(), err)
}
