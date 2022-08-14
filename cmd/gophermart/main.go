package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/server"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
)

func main() {
	panic(1)
	log.SetOutput(os.Stdout)
	cfg := config.New()

	cfg.Init()
	flag.Parse()

	err := storage.RunMigrations(cfg.DatabaseURI)
	if err != nil {
		panic(fmt.Sprintf("could not run migrations: %v", err))
	}

	userRepo, err := storage.NewUserRepository(cfg.DatabaseURI)
	if err != nil {
		panic(fmt.Sprintf("could not initialize user repo: %v", err))
	}

	ordersRepo, err := storage.NewOrdersRepository(cfg.DatabaseURI)
	if err != nil {
		panic(fmt.Sprintf("could not initialize orders repo: %v", err))
	}

	balanceRepo, err := storage.NewBalanceRepository(cfg.DatabaseURI)
	if err != nil {
		panic(fmt.Sprintf("could not initialize balance repo: %v", err))
	}

	maxRequestsPerSecond := 50 // todo: move max requests per sec to config
	accrualService := services.NewAccrualHTTPClient(http.DefaultClient, cfg.AccrualSystemAddress, maxRequestsPerSecond)
	auth := services.NewAuth(userRepo, cfg.JWTSecret)
	balanceProcessor := services.NewBalanceProcessor(balanceRepo)
	ordersProcessor := services.NewOrdersProcessor(ordersRepo, balanceProcessor, accrualService)
	srv := server.New(cfg, auth, ordersProcessor, balanceProcessor)

	srv.Run()
}
