package storage

import (
	"errors"
	"fmt"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"os"

	"github.com/golang-migrate/migrate/v4"
)

type OrdersStorage interface {
	IsRegistered(orderID int) (bool, error)
	RegisterOrder(orderID int, items []models.OrderItem) error
}

func RunMigrations(dsn string) error {
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
