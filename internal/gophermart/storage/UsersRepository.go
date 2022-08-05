package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

type Users interface {
	CreateNew(login string, password string) (models.User, error)
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

	if err = runMigrations(dsn, "./migrations"); err != nil {
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

func (u *UsersRepository) CreateNew(login string, password string) (models.User, error) {
	// TODO: implement
	panic("implement me")
}
