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

const testDSN = "postgres://postgres:postgres@db:5432/praktikum?sslmode=disable"

type BalanceRepositoryTestSuite struct {
	suite.Suite
	ordersRepository  *OrdersRepository
	usersRepository   *UsersRepository
	balanceRepository *BalanceRepository
}

func TestBalanceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceRepositoryTestSuite))
}

func (s *BalanceRepositoryTestSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = testDSN
	}
	ordersRepository, err := NewOrdersRepository(dsn)
	require.NoError(s.T(), err)
	s.ordersRepository = ordersRepository

	usersRepository, err := NewUserRepository(dsn)
	require.NoError(s.T(), err)
	s.usersRepository = usersRepository

	balanceRepository, err := NewBalanceRepository(dsn)
	require.NoError(s.T(), err)
	s.balanceRepository = balanceRepository

	err = RunMigrations(dsn)
	require.NoError(s.T(), err)
}

func (s *BalanceRepositoryTestSuite) TearDownTest() {
	_, _ = s.ordersRepository.conn.Exec(context.Background(), "truncate table orders cascade")
	_, _ = s.ordersRepository.conn.Exec(context.Background(), "truncate table users cascade")
	_, _ = s.ordersRepository.conn.Exec(context.Background(), "truncate table withdraws cascade")
}

func (s *BalanceRepositoryTestSuite) exists(orderID int, userID int, amount float64) bool {
	var exists bool

	err := s.ordersRepository.conn.QueryRow(
		context.Background(),
		"select exists(select 1 from withdraws where order_id = $1 and user_id = $2 and amount = $3)",
		orderID,
		userID,
		amount,
	).Scan(&exists)

	assert.NoError(s.T(), err)
	return exists
}

func (s *BalanceRepositoryTestSuite) TestAddWithdraw() {
	user, err := s.usersRepository.CreateNew("login", "password")
	require.NoError(s.T(), err)

	withdrawAmount := 100.0
	orderID := 100

	err = s.balanceRepository.AddWithdraw(orderID, user.ID, withdrawAmount)
	assert.NoError(s.T(), err)
	assert.True(s.T(), s.exists(orderID, user.ID, withdrawAmount))
}

func (s *BalanceRepositoryTestSuite) TestGetTotalWithdraws() {
	max, err := s.usersRepository.CreateNew("max", "password")
	require.NoError(s.T(), err)
	john, err := s.usersRepository.CreateNew("john", "password")
	require.NoError(s.T(), err)

	_ = s.balanceRepository.AddWithdraw(10, max.ID, 10.0)
	_ = s.balanceRepository.AddWithdraw(20, max.ID, 20.0)
	_ = s.balanceRepository.AddWithdraw(30, john.ID, 40.0)

	totalWithdraws, err := s.balanceRepository.GetTotalWithdraws(max.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 30.0, totalWithdraws)
}

func (s *BalanceRepositoryTestSuite) TestGetTotalAccrual() {
	max, err := s.usersRepository.CreateNew("max", "password")
	require.NoError(s.T(), err)
	john, err := s.usersRepository.CreateNew("john", "password")
	require.NoError(s.T(), err)

	s.addAccrual(10, max.ID, 10.0)
	s.addAccrual(20, max.ID, 20.0)
	s.addAccrual(40, john.ID, 40.0)

	totalAccrual, err := s.balanceRepository.GetTotalAccrual(max.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 30.0, totalAccrual)
}

func (s *BalanceRepositoryTestSuite) addAccrual(orderID int, userID int, amount float64) {
	_, err := s.balanceRepository.conn.Exec(
		context.Background(),
		"insert into orders (id, created_by, status, accrual) values ($1, $2, $3, $4)",
		orderID,
		userID,
		models.OrderStatusProcessed,
		amount,
	)
	require.NoError(s.T(), err)
}
