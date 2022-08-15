package main

import (
	"flag"

	"github.com/belamov/ypgo-gophermart/internal/accrual/config"
	"github.com/belamov/ypgo-gophermart/internal/accrual/server"
	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.New()

	cfg.Init()
	flag.Parse()

	rewardsRepo, err := storage.NewRewardsRepository(cfg.DatabaseURI)
	if err != nil {
		panic(err.Error())
	}

	ordersRepo, err := storage.NewOrdersRepository(cfg.DatabaseURI)
	if err != nil {
		panic(err.Error())
	}

	ordersManager := services.NewOrderManager(ordersRepo)

	srv := server.New(cfg, ordersManager, rewardsRepo)

	log.Fatal().Err(srv.Run())
}
