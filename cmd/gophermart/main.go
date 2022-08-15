package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/server"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	cfg := config.New()

	cfg.Init()
	flag.Parse()

	log.Info().Msgf("config initialized: %+v", cfg)

	err := storage.RunMigrations(cfg.DatabaseURI)
	if err != nil {
		log.Panic().Err(err).Msg("could not run migrations")
	}

	userRepo, err := storage.NewUserRepository(cfg.DatabaseURI)
	if err != nil {
		log.Panic().Err(err).Msg("could not initialize user repo")
	}

	ordersRepo, err := storage.NewOrdersRepository(cfg.DatabaseURI)
	if err != nil {
		log.Panic().Err(err).Msg("could not initialize orders repo")
	}

	balanceRepo, err := storage.NewBalanceRepository(cfg.DatabaseURI)
	if err != nil {
		log.Panic().Err(err).Msg("could not initialize balance repo")
	}

	accrualService := services.NewAccrualHTTPClient(http.DefaultClient, cfg.AccrualSystemAddress, cfg.MaxRequestsPerSecondsToAccrual)
	auth := services.NewAuth(userRepo, cfg.JWTSecret)
	balanceProcessor := services.NewBalanceProcessor(balanceRepo)
	ordersProcessor := services.NewOrdersProcessor(ordersRepo, balanceProcessor, accrualService)
	srv := server.New(cfg, auth, ordersProcessor, balanceProcessor)

	srv.Run()
}
