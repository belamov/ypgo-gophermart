package server

import (
	"log"
	"net/http"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/handlers"
)

type Server struct {
	config *config.Config
}

func (s *Server) Run() {
	r := handlers.NewRouter()

	httpServer := &http.Server{
		Addr:              s.config.RunAddress,
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func New(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}
