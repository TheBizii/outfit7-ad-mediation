package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	PSQLUrl string
	AppPort string
}

func Load() *Config {
	// load the local .env file
	_ = godotenv.Load()

	cfg := &Config{
		PSQLUrl: os.Getenv("PSQL_URL"),
		AppPort: os.Getenv("APP_PORT"),
	}

	// make sure the user is properly made aware of missing ENV variables
	if cfg.PSQLUrl == "" {
		log.Fatal("Missing PSQL_URL from .env")
	}
	if cfg.AppPort == "" {
		log.Fatal("Missing APP_PORT from .env")
	}

	return cfg
}
