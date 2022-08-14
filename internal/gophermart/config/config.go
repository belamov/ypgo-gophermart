package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecret            string
}

func New() *Config {
	return &Config{
		RunAddress:           "http://localhost:8080",
		DatabaseURI:          "",
		AccrualSystemAddress: "",
		JWTSecret:            "secret",
	}
}

func (c *Config) Init() {
	flag.StringVar(&c.RunAddress, "a", getEnv("RUN_ADDRESS", ":8080"), "host to listen on")
	flag.StringVar(&c.AccrualSystemAddress, "f", getEnv("ACCRUAL_SYSTEM_ADDRESS", ""), "file storage path")
	flag.StringVar(&c.DatabaseURI, "d", getEnv("DATABASE_URI", ""), "database dsn for connecting to postgres")
	flag.StringVar(&c.JWTSecret, "js", getEnv("JWT_SECRET", "some jwt secret token"), "secret token for signing and verifying jwt tokens")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Default().Println("found value in env: " + value)
		return value
	}
	return fallback
}
