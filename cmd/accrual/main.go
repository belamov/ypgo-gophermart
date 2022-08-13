package main

import (
	"flag"

	"github.com/belamov/ypgo-gophermart/internal/accrual/config"
	"github.com/belamov/ypgo-gophermart/internal/accrual/server"
)

func main() {
	cfg := config.New()

	cfg.Init()
	flag.Parse()

	srv := server.New(cfg)

	srv.Run()
}
