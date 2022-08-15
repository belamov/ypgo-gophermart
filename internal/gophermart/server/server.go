package server

import (
	"context"
	"net/http"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/handlers"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
)

type Server struct {
	server *http.Server
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func New(
	config *config.Config,
	auth services.Auth,
	ordersManager services.OrdersManagerInterface,
	balanceProcessor services.BalanceProcessorInterface,
) *Server {
	r := handlers.NewRouter(auth, ordersManager, balanceProcessor)
	httpServer := &http.Server{
		Addr:              config.RunAddress,
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
	}
	return &Server{
		server: httpServer,
	}
}
