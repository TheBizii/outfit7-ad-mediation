package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err := db.Ping(); err != nil {
        	log.Fatal("Failed to ping DB:", err)
    	}

	log.Println("Connected to Postgres!")
	return db
}
