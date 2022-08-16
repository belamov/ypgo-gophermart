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
	err := s.ordersRepository.CreateNew(orderID, orderItems)
	assert.NoError(s.T(), err)

	err = s.ordersRepository.CreateNew(orderID, orderItems)
	var nue *NotUniqueError
	assert.ErrorAs(s.T(), err, &nue)
}
