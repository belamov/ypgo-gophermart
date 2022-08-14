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

func TestUsersRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UsersRepositoryTestSuite))
}

func (s *UsersRepositoryTestSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db_gophermart:5432/praktikum?sslmode=disable"
	}
	repo, err := NewUserRepository(dsn)
	require.NoError(s.T(), err)
	s.repo = repo

	err = RunMigrations(dsn)
	require.NoError(s.T(), err)
}

func (s *UsersRepositoryTestSuite) TearDownTest() {
	conn, _ := s.repo.pool.Acquire(context.Background())
	_, err := conn.Exec(context.Background(), "truncate table users cascade")
	conn.Release()
	require.NoError(s.T(), err)
}

func (s *UsersRepositoryTestSuite) exists(user models.User) bool {
	var exists bool
	conn, _ := s.repo.pool.Acquire(context.Background())
	err := conn.QueryRow(
		context.Background(),
		"select exists(select 1 from users where login = $1 and password = $2)",
		user.Login,
		user.HashedPassword,
	).Scan(&exists)
	conn.Release()
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

func (s *UsersRepositoryTestSuite) TestFindByLogin() {
	credentials := models.Credentials{
		Login:    "login",
		Password: "hash",
	}
	existingUser, err := s.repo.CreateNew(credentials.Login, credentials.Password)
	require.NoError(s.T(), err)
	assert.True(s.T(), s.exists(existingUser))

	fetchedUser, err := s.repo.FindByLogin(credentials.Login)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), existingUser.ID, fetchedUser.ID)
	assert.Equal(s.T(), existingUser.Login, fetchedUser.Login)
	assert.Equal(s.T(), existingUser.HashedPassword, fetchedUser.HashedPassword)

	nonExistingLogin := "non existing"
	notFoundUser, err := s.repo.FindByLogin(nonExistingLogin)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), models.User{}, notFoundUser)
}
