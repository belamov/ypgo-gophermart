package server

import (
	"log"
	"net/http"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/accrual/config"
	"github.com/belamov/ypgo-gophermart/internal/accrual/handlers"
	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
)

type Server struct {
	config       *config.Config
	orderManager services.OrderManagementInterface
}

func (s *Server) Run() {
	r := handlers.NewRouter(s.orderManager)

	httpServer := &http.Server{
		Addr:              s.config.RunAddress,
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func New(
	config *config.Config,
) *Server {
	return &Server{
		config: config,
	}
}
