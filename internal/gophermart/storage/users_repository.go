package storage

import (
	"context"
	"errors"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type UsersRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(dsn string) (*UsersRepository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &UsersRepository{
		conn: conn,
	}, nil
}

func (repo *UsersRepository) CreateNew(login string, password string) (models.User, error) {
	user := models.User{
		Login:          login,
		HashedPassword: password,
	}

	err := repo.conn.QueryRow(
		context.Background(),
		"insert into users (login, password) values ($1, $2) returning id",
		user.Login,
		user.HashedPassword,
	).Scan(&user.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return models.User{}, NewNotUniqueError("login", err)
		}
	}
	return user, err
}

func (repo *UsersRepository) FindByLogin(login string) (models.User, error) {
	var user models.User
	err := repo.conn.QueryRow(
		context.Background(),
		"select id, login, password from users where login=$1",
		login,
	).Scan(&user.ID, &user.Login, &user.HashedPassword)

	if err == pgx.ErrNoRows {
		return models.User{}, nil
	}

	return user, err
}
