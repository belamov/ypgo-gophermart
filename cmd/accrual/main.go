package main

import (
	"flag"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/config"
)

func main() {
	cfg := config.New()

	cfg.Init()
	flag.Parse()
}
