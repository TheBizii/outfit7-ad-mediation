package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PSQLUrl string
	AppPort string
}

func Load() (*Config, error) {
	// load the local .env file
	_ = godotenv.Load()

	// define the Postgres connection
	PSQL_HOST := os.Getenv("PSQL_HOST")
	PSQL_USER := os.Getenv("PSQL_USER")
	PSQL_PASSWORD := os.Getenv("PSQL_PASSWORD")
	PSQL_DBNAME := os.Getenv("PSQL_DBNAME")
	PSQL_PORT := os.Getenv("PSQL_PORT")

	psqlUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		PSQL_HOST,
		PSQL_USER,
		PSQL_PASSWORD,
		PSQL_DBNAME,
		PSQL_PORT)

	// init the config struct
	cfg := &Config{
		PSQLUrl: psqlUrl,
		AppPort: os.Getenv("APP_PORT"),
	}

	// make sure the user is properly made aware of missing ENV variables
	if PSQL_HOST == "" || PSQL_USER == "" || PSQL_PASSWORD == "" || PSQL_DBNAME == "" || PSQL_PORT == "" {
		return nil, fmt.Errorf("Make sure your .env config has the following required variables present: PSQL_HOST, PSQL_USER, PSQL_PASSWORD, PSQL_DBNAME, PSQL_PORT.")
	}
	if cfg.AppPort == "" {
		return nil, fmt.Errorf("Missing APP_PORT from .env.")
	}

	return cfg, nil
}
