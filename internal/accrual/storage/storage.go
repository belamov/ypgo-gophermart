package storage

import (
	"errors"
	"os"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

type OrdersStorage interface {
	CreateNew(orderID int, items []models.OrderItem) (models.Order, error)
	AddAccrual(orderID int, accrual float64) error
	ChangeStatus(orderID int, status models.OrderStatus) error
	GetOrdersForProcessing() ([]models.Order, error)
	GetOrder(orderID int) (models.Order, error)
}

type RewardsStorage interface {
	Save(rewardCondition models.Reward) error
	GetMatchingReward(orderItem models.OrderItem) (models.Reward, error)
}

func RunMigrations(dsn string) error {
	m, err := migrate.New(getMigrationsPath(), dsn)
	if err != nil {
		return err
	}

	log.Info().Msg("Migrating...")

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info().Msg("Nothing to migrate")
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("Migration failed!")
		return err
	}

	log.Info().Msg("Migrated successfully")
	return nil
}

func getMigrationsPath() string {
	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = "file://./migrations"
	}
	return path
}
