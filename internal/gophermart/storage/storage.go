package storage

import (
	"errors"
	"fmt"
	"os"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang-migrate/migrate/v4"
)

type OrdersStorage interface {
	CreateNew(orderID int, userID int) (models.Order, error)
	FindByID(orderID int) (models.Order, error)
	GetUsersOrders(userID int) ([]models.Order, error)
}

type UsersStorage interface {
	CreateNew(login string, password string) (models.User, error)
	FindByLogin(login string) (models.User, error)
}

type BalanceStorage interface {
	AddWithdraw(orderID int, userID int, withdrawAmount float64) error
	GetTotalAccrual(userID int) (float64, error)
	GetTotalWithdraws(userID int) (float64, error)
}

func runMigrations(dsn string) error {
	m, err := migrate.New(getMigrationsPath(), dsn)
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

func getMigrationsPath() string {
	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = "file://./migrations"
	}
	return path
}
