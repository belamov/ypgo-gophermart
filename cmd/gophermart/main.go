package main

import (
	"flag"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/server"
)

func main() {
	cfg := config.New()
	auth := &services.Auth{}

	cfg.Init()
	flag.Parse()

	srv := server.New(cfg, auth)

	srv.Run()
}
