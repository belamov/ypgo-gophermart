package storage

import (
	"context"
	"errors"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type UsersRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(dsn string) (*UsersRepository, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &UsersRepository{
		pool: pool,
	}, nil
}

func (repo *UsersRepository) CreateNew(login string, password string) (models.User, error) {
	user := models.User{
		Login:          login,
		HashedPassword: password,
	}

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldnt acquire connection from pool")
		return models.User{}, err
	}

	err = conn.QueryRow(
		context.Background(),
		"insert into users (login, password) values ($1, $2) returning id",
		user.Login,
		user.HashedPassword,
	).Scan(&user.ID)

	conn.Release()

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

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldnt acquire connection from pool")
		return user, err
	}

	err = conn.QueryRow(
		context.Background(),
		"select id, login, password from users where login=$1",
		login,
	).Scan(&user.ID, &user.Login, &user.HashedPassword)

	conn.Release()

	if err == pgx.ErrNoRows {
		return models.User{}, nil
	}

	return user, err
}
