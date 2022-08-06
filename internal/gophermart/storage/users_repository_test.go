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

type UsersRepositoryTestSuite struct {
	suite.Suite
	repo *UsersRepository
}

func TestPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UsersRepositoryTestSuite))
}

func (s *UsersRepositoryTestSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/praktikum?sslmode=disable"
	}
	repo, err := NewUserRepository(dsn)
	require.NoError(s.T(), err)
	s.repo = repo
}

func (s *UsersRepositoryTestSuite) TearDownTest() {
	_, err := s.repo.conn.Exec(context.Background(), "truncate table users")
	require.NoError(s.T(), err)
}

func (s *UsersRepositoryTestSuite) exists(user models.User) bool {
	var exists bool

	err := s.repo.conn.QueryRow(
		context.Background(),
		"select exists(select 1 from users where login = $1 and password = $2)",
		user.Login,
		user.HashedPassword,
	).Scan(&exists)

	assert.NoError(s.T(), err)
	return exists
}

func (s *UsersRepositoryTestSuite) TestCreate() {
	login := "login"
	passwordHash := "password"
	createdUser, err := s.repo.CreateNew(login, passwordHash)
	require.NoError(s.T(), err)
	assert.True(s.T(), s.exists(createdUser))
}

func (s *UsersRepositoryTestSuite) TestUniqueLogin() {
	login := "login"
	passwordHash := "password"
	existingUser, err := s.repo.CreateNew(login, passwordHash)
	require.NoError(s.T(), err)
	assert.True(s.T(), s.exists(existingUser))

	createdUser, err := s.repo.CreateNew(login, passwordHash)
	var notUniqueError *NotUniqueError
	assert.ErrorAs(s.T(), err, &notUniqueError)
	assert.Equal(s.T(), "login", notUniqueError.Field)
	assert.False(s.T(), s.exists(createdUser))
}
