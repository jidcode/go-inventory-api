package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ventry/internal/domain/config"
)

var DB *sqlx.DB

func Connect(config config.Variables) *sqlx.DB {
	connectionString := config.DatabaseUrl

	DB, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	log.Print("Database connection successful!")

	return DB
}
