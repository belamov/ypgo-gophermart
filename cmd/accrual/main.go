package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/accrual/config"
	"github.com/belamov/ypgo-gophermart/internal/accrual/server"
	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	cfg := initConfig()
	rewardsRepo, ordersRepo := initRepos(cfg)
	accrualProcessor := services.NewAccrualProcessor(ordersRepo, rewardsRepo)
	ordersManager := services.NewOrderManager(ordersRepo, accrualProcessor)
	srv := server.New(cfg, ordersManager, rewardsRepo)

	ctx, cancel := context.WithCancel(context.Background())

	go accrualProcessor.StartProcessing(ctx)

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
	time.Sleep(time.Second * 3) //nolint:gomnd

	log.Info().Msg("goodbye")
}

func initConfig() *config.Config {
	cfg := config.New()

	cfg.Init()
	flag.Parse()

	log.Info().Msgf("config initialized: %+v", cfg)
	return cfg
}

func initRepos(cfg *config.Config) (*storage.RewardsRepository, *storage.OrdersRepository) {
	err := storage.RunMigrations(cfg.DatabaseURI)
	if err != nil {
		log.Panic().Err(err).Msg("could not run migrations")
	}

	rewardsRepo, err := storage.NewRewardsRepository(cfg.DatabaseURI)
	if err != nil {
		log.Panic().Err(err).Msg("could not initialize rewards repo")
	}

	ordersRepo, err := storage.NewOrdersRepository(cfg.DatabaseURI)
	if err != nil {
		log.Panic().Err(err).Msg("could not initialize orders repo")
	}

	return rewardsRepo, ordersRepo
}
