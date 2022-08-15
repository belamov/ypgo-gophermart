package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/server"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	cfg := initConfig()
	userRepo, ordersRepo, balanceRepo := initRepos(cfg)

	accrualService := services.NewAccrualHTTPClient(http.DefaultClient, cfg.AccrualSystemAddress, cfg.MaxRequestsPerSecondsToAccrual)
	balanceProcessor := services.NewBalanceProcessor(balanceRepo)

	ordersProcessor := services.NewOrderProcessor(ordersRepo, accrualService, balanceProcessor)

	ordersManager := services.NewOrdersManager(ordersRepo, balanceProcessor, ordersProcessor)
	auth := services.NewAuth(userRepo, cfg.JWTSecret)
	srv := server.New(cfg, auth, ordersManager, balanceProcessor)

	ctx, cancel := context.WithCancel(context.Background())

	go ordersProcessor.StartProcessing(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server couldn't start!")
		}
	}()
	log.Info().Msg("Server Started")

	<-done
	log.Info().Msg("Shutting down gracefully")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server Shutdown Failed")
	}
	log.Info().Msg("Server shut down successfully")

	log.Info().Msg("canceling processing orders...")

	cancel()

	log.Info().Msg("goodbye")
}

func initConfig() *config.Config {
	cfg := config.New()

	cfg.Init()
	flag.Parse()

	log.Info().Msgf("config initialized: %+v", cfg)
	return cfg
}

func initRepos(cfg *config.Config) (*storage.UsersRepository, *storage.OrdersRepository, *storage.BalanceRepository) {
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

	return userRepo, ordersRepo, balanceRepo
}
