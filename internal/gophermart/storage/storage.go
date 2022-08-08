package storage

import (
	"errors"
	"fmt"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang-migrate/migrate/v4"
)

type OrdersStorage interface {
	CreateNew(orderID int, userID int) (models.Order, error)
	FindByID(orderID int) (models.Order, error)
}

type UsersStorage interface {
	CreateNew(login string, password string) (models.User, error)
	FindByLogin(login string) (models.User, error)
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
