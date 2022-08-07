package server

import (
	"log"
	"net/http"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/handlers"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
)

type Server struct {
	config *config.Config
	auth   services.Auth
	orders services.OrdersProcessorInterface
}

func (s *Server) Run() {
	r := handlers.NewRouter(s.auth, s.orders)

	httpServer := &http.Server{
		Addr:              s.config.RunAddress,
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func New(config *config.Config, auth services.Auth, orders services.OrdersProcessorInterface) *Server {
	return &Server{
		config: config,
		auth:   auth,
		orders: orders,
	}
}
