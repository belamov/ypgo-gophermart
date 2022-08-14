package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const testDSN = "postgres://postgres:postgres@db_gophermart:5432/praktikum?sslmode=disable"

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
	conn, _ := s.balanceRepository.pool.Acquire(context.Background())
	_, _ = conn.Exec(context.Background(), "truncate table orders cascade")
	_, _ = conn.Exec(context.Background(), "truncate table users cascade")
	_, _ = conn.Exec(context.Background(), "truncate table withdraws cascade")
	conn.Release()
}

func (s *BalanceRepositoryTestSuite) exists(orderID int, userID int, amount float64) bool {
	var exists bool
	conn, _ := s.balanceRepository.pool.Acquire(context.Background())

	err := conn.QueryRow(
		context.Background(),
		"select exists(select 1 from withdraws where order_id = $1 and user_id = $2 and amount = $3)",
		orderID,
		userID,
		amount,
	).Scan(&exists)

	conn.Release()
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

	totalWithdraws, err := s.balanceRepository.GetTotalWithdrawAmount(max.ID)
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

	totalAccrual, err := s.balanceRepository.GetTotalAccrualAmount(max.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 30.0, totalAccrual)
}

func (s *BalanceRepositoryTestSuite) addAccrual(orderID int, userID int, amount float64) {
	conn, _ := s.balanceRepository.pool.Acquire(context.Background())

	_, err := conn.Exec(
		context.Background(),
		"insert into orders (id, created_by, status, accrual) values ($1, $2, $3, $4)",
		orderID,
		userID,
		models.OrderStatusProcessed,
		amount,
	)

	conn.Release()

	require.NoError(s.T(), err)
}

func (s *BalanceRepositoryTestSuite) TestGetUserWithdrawals() {
	max, err := s.usersRepository.CreateNew("max", "password")
	require.NoError(s.T(), err)
	john, err := s.usersRepository.CreateNew("john", "password")
	require.NoError(s.T(), err)

	oldMaxWithdrawalOrderID := 123
	newMaxWithdrawalOrderID := 456
	err = s.balanceRepository.AddWithdraw(oldMaxWithdrawalOrderID, max.ID, 10)
	require.NoError(s.T(), err)
	err = s.balanceRepository.AddWithdraw(999, john.ID, 10)
	require.NoError(s.T(), err)
	time.Sleep(time.Second)
	err = s.balanceRepository.AddWithdraw(newMaxWithdrawalOrderID, max.ID, 20)
	require.NoError(s.T(), err)

	maxesWithdrawals, err := s.balanceRepository.GetUserWithdrawals(max.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), maxesWithdrawals, 2)
	assert.Equal(s.T(), oldMaxWithdrawalOrderID, maxesWithdrawals[0].OrderID)
	assert.Equal(s.T(), newMaxWithdrawalOrderID, maxesWithdrawals[1].OrderID)
}

func (s *BalanceRepositoryTestSuite) TestAddAccrual() {
	user, err := s.usersRepository.CreateNew("max", "password")
	require.NoError(s.T(), err)

	order, err := s.ordersRepository.CreateNew(1, user.ID)
	require.NoError(s.T(), err)

	accrualAmount := 10.5

	newOrder, err := s.ordersRepository.FindByID(order.ID)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 0.0, newOrder.Accrual)
	assert.Equal(s.T(), models.OrderStatusNew, newOrder.Status)

	err = s.balanceRepository.AddAccrual(order.ID, accrualAmount)
	assert.NoError(s.T(), err)

	updatedOrder, err := s.ordersRepository.FindByID(order.ID)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), accrualAmount, updatedOrder.Accrual)
	assert.Equal(s.T(), models.OrderStatusProcessed, updatedOrder.Status)
}
