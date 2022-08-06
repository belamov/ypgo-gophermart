package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type Users interface {
	CreateNew(login string, password string) (models.User, error)
	FindByLogin(login string) (models.User, error)
}

type UsersRepository struct {
	Dsn  string
	conn *pgx.Conn
}

func NewUserRepository(dsn string) (*UsersRepository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = runMigrations(dsn, "file://./migrations"); err != nil {
		return nil, err
	}

	return &UsersRepository{
		Dsn:  dsn,
		conn: conn,
	}, nil
}

func runMigrations(dsn string, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Nothing to migrate")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Println("Migrated successfully")
	return nil
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
