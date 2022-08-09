package storage

import (
	"context"
	"os"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrdersRepositoryTestSuite struct {
	suite.Suite
	ordersRepository *OrdersRepository
	usersRepository  *UsersRepository
}

func TestOrdersRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersRepositoryTestSuite))
}

func (s *OrdersRepositoryTestSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/praktikum?sslmode=disable"
	}
	ordersRepository, err := NewOrdersRepository(dsn)
	require.NoError(s.T(), err)
	s.ordersRepository = ordersRepository

	usersRepository, err := NewUserRepository(dsn)
	require.NoError(s.T(), err)
	s.usersRepository = usersRepository
}

func (s *OrdersRepositoryTestSuite) TearDownTest() {
	_, _ = s.ordersRepository.conn.Exec(context.Background(), "truncate table orders cascade")
	_, _ = s.ordersRepository.conn.Exec(context.Background(), "truncate table users cascade")
}

func (s *OrdersRepositoryTestSuite) exists(order models.Order) bool {
	var exists bool

	err := s.ordersRepository.conn.QueryRow(
		context.Background(),
		"select exists(select 1 from orders where id = $1 and created_by = $2)",
		order.ID,
		order.CreatedBy,
	).Scan(&exists)

	assert.NoError(s.T(), err)
	return exists
}

func (s *OrdersRepositoryTestSuite) TestCreateNew() {
	user, err := s.usersRepository.CreateNew("login", "password")
	require.NoError(s.T(), err)

	orderID := 123

	createdOrder, err := s.ordersRepository.CreateNew(orderID, user.ID)
	require.NoError(s.T(), err)
	expectedCreatedOrder := models.Order{ID: orderID, CreatedBy: user.ID}
	assert.True(s.T(), s.exists(expectedCreatedOrder))
	assert.Equal(s.T(), orderID, createdOrder.ID)
	assert.Equal(s.T(), user.ID, createdOrder.CreatedBy)
	assert.Equal(s.T(), models.OrderStatusNew, createdOrder.Status)
}

func (s *OrdersRepositoryTestSuite) TestFindByID() {
	user, err := s.usersRepository.CreateNew("login", "password")
	require.NoError(s.T(), err)

	orderID := 123

	createdOrder, err := s.ordersRepository.CreateNew(orderID, user.ID)
	require.NoError(s.T(), err)

	fetchedOrder, err := s.ordersRepository.FindByID(createdOrder.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), createdOrder.ID, fetchedOrder.ID)
	assert.Equal(s.T(), createdOrder.CreatedBy, fetchedOrder.CreatedBy)
	assert.Equal(s.T(), createdOrder.Status, fetchedOrder.Status)

	notFoundOrder, err := s.ordersRepository.FindByID(111111)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, notFoundOrder.ID)
	assert.Equal(s.T(), 0, notFoundOrder.CreatedBy)
}
