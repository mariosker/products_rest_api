package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL      string
	ServerHost string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	cfg := &Config{
		DBURL:      getEnv("DBURL", "postgres://postgres:postgres@postgres/products"),
		ServerHost: getEnv("ServerHost", "localhost"),
		ServerPort: getEnv("ServerPort", "8080"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
