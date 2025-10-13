package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Connect(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Failed to ping the database: %w", err)
	}

	log.Println("Connected to Postgres!")
	Conn = db
	return nil
}
