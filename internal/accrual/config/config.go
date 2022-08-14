package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress  string
	DatabaseURI string
}

func New() *Config {
	return &Config{
		RunAddress:  "http://localhost:8081",
		DatabaseURI: "",
	}
}

func (c *Config) Init() {
	flag.StringVar(&c.RunAddress, "a", getEnv("RUN_ADDRESS", ":8081"), "host to listen on")
	flag.StringVar(&c.DatabaseURI, "d", getEnv("DATABASE_URI", ""), "database dsn for connecting to postgres")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
