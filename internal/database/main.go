package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type DBConnection = *pgx.Conn

var db *pgx.Conn

// InitDB initializes the database connection
func InitDB(connectionString string) error {
	var err error
	db, err = pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return err
	}
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if err := db.Close(context.Background()); err != nil {
		log.Fatal("Error closing database connection:", err)
	}
}

// GetDB returns the database connection
func GetDB() *pgx.Conn {
	return db
}
