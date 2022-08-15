package server

import (
	"net/http"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/accrual/config"
	"github.com/belamov/ypgo-gophermart/internal/accrual/handlers"
	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
)

type Server struct {
	server *http.Server
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func New(
	config *config.Config,
	orderManager services.OrderManagementInterface,
	rewardsStorage storage.RewardsStorage,
) *Server {
	r := handlers.NewRouter(orderManager, rewardsStorage)

	httpServer := &http.Server{
		Addr:              config.RunAddress,
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
	}
	return &Server{server: httpServer}
}
