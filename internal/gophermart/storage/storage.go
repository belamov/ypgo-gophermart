package storage

import (
	"errors"
	"fmt"
	"os"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type OrdersStorage interface {
	CreateNew(orderID int, userID int) (models.Order, error)
	FindByID(orderID int) (models.Order, error)
	GetUsersOrders(userID int) ([]models.Order, error)
	ChangeStatus(order models.Order, status models.OrderStatus) error
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
		path = "file://internal/gophermart/storage/migrations"
	}
	return path
}
