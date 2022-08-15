package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	RunAddress                     string
	DatabaseURI                    string
	AccrualSystemAddress           string
	JWTSecret                      string
	MaxRequestsPerSecondsToAccrual int
}

func New() *Config {
	defaultMaxRPS := 50
	return &Config{
		RunAddress:                     "http://localhost:8080",
		DatabaseURI:                    "",
		AccrualSystemAddress:           "",
		JWTSecret:                      "secret",
		MaxRequestsPerSecondsToAccrual: defaultMaxRPS,
	}
}

func (c *Config) Init() {
	flag.StringVar(&c.RunAddress, "a", getEnv("RUN_ADDRESS", ":8080"), "host to listen on")
	flag.StringVar(&c.AccrualSystemAddress, "f", getEnv("ACCRUAL_SYSTEM_ADDRESS", ""), "file storage path")
	flag.StringVar(&c.DatabaseURI, "d", getEnv("DATABASE_URI", ""), "database dsn for connecting to postgres")
	flag.StringVar(&c.JWTSecret, "js", getEnv("JWT_SECRET", "some jwt secret token"), "secret token for signing and verifying jwt tokens")
	maxRequestsFromEnv, err := strconv.Atoi(getEnv("MAX_REQUESTS_PER_SECOND_TO_ACCRUAL", "50"))
	if err != nil {
		log.Panic(err)
	}
	flag.IntVar(&c.MaxRequestsPerSecondsToAccrual, "mrps", maxRequestsFromEnv, "maximum requests per seconds to accrual service allowed. used for throttling outcoming requests")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Println("found value in env: " + value)
		return value
	}
	return fallback
}
