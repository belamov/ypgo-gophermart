package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/handlers"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
)

type Server struct {
	config           *config.Config
	auth             services.Auth
	ordersProcessor  services.OrdersProcessorInterface
	balanceProcessor services.BalanceProcessorInterface
}

func (s *Server) Run() {
	r := handlers.NewRouter(s.auth, s.ordersProcessor, s.balanceProcessor)
	fmt.Println("run address: " + s.config.RunAddress)
	fmt.Println("accrual address: " + s.config.AccrualSystemAddress)
	httpServer := &http.Server{
		Addr:              s.config.RunAddress,
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func New(
	config *config.Config,
	auth services.Auth,
	ordersProcessor services.OrdersProcessorInterface,
	balanceProcessor services.BalanceProcessorInterface,
) *Server {
	return &Server{
		config:           config,
		auth:             auth,
		ordersProcessor:  ordersProcessor,
		balanceProcessor: balanceProcessor,
	}
}
