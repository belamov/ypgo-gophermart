package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func New() *Config {
	return &Config{
		RunAddress:           "http://localhost:8080",
		DatabaseURI:          "",
		AccrualSystemAddress: "",
	}
}

func (c *Config) Init() {
	flag.StringVar(&c.RunAddress, "a", getEnv("RUN_ADDRESS", ":8080"), "host to listen on")
	flag.StringVar(&c.AccrualSystemAddress, "f", getEnv("ACCRUAL_SYSTEM_ADDRESS", ""), "file storage path")
	flag.StringVar(&c.DatabaseURI, "d", getEnv("DATABASE_URI", ""), "database dsn for connecting to postgres")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
