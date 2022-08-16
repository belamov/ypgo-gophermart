package storage

import (
	"errors"
	"os"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

type OrdersStorage interface {
	CreateNew(orderID int, userID int) (models.Order, error)
	FindByID(orderID int) (models.Order, error)
	GetUsersOrders(userID int) ([]models.Order, error)
	ChangeStatus(order models.Order, status models.OrderStatus) error
	GetOrdersForProcessing() ([]models.Order, error)
}

type UsersStorage interface {
	CreateNew(login string, password string) (models.User, error)
	FindByLogin(login string) (models.User, error)
}

type BalanceStorage interface {
	AddWithdraw(orderID int, userID int, withdrawAmount float64) error
	GetTotalAccrualAmount(userID int) (float64, error)
	GetTotalWithdrawAmount(userID int) (float64, error)
	GetUserWithdrawals(userID int) ([]models.Withdrawal, error)
	AddAccrual(orderID int, accrual float64) error
}

func RunMigrations(dsn string) error {
	m, err := migrate.New(getMigrationsPath(), dsn+"&x-migrations-table=accrual_migrations")
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
		path = "file://internal/gophermart/storage/migrations"
	}
	return path
}
