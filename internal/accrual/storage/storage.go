package storage

import (
	"errors"
	"fmt"
	"os"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type OrdersStorage interface {
	CreateNew(orderID int, items []models.OrderItem) error
}

type RewardsStorage interface {
	Save(rewardCondition models.RewardCondition) error
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
