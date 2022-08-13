package main

import (
	"flag"

	"github.com/belamov/ypgo-gophermart/internal/accrual/config"
)

func main() {
	cfg := config.New()

	cfg.Init()
	flag.Parse()
}
