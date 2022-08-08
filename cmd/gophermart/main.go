package main

import (
	"flag"
	"fmt"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/server"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
)

func main() {
	cfg := config.New()

	cfg.Init()
	flag.Parse()

	userRepo, err := storage.NewUserRepository(cfg.DatabaseURI)
	if err != nil {
		panic(fmt.Sprintf("could not initialize user repo: %v", err))
	}

	ordersRepo, err := storage.NewOrdersRepository(cfg.DatabaseURI)
	if err != nil {
		panic(fmt.Sprintf("could not initialize orders repo: %v", err))
	}

	auth := services.NewAuth(userRepo, cfg.JWTSecret)
	orders := services.NewOrdersProcessor(ordersRepo)
	srv := server.New(cfg, auth, orders)

	srv.Run()
}
