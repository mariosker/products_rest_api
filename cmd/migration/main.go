package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mariosker/products_rest_api/internal/config"
)

func main() {
	// Parse command-line arguments for migration direction
	var direction string
	flag.StringVar(&direction, "direction", "up", "Specify 'up' or 'down' for migration")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize pgx configuration
	pgxConfig, err := pgx.ParseConfig(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to parse pgx config: %v", err)
	}

	// Use pgx stdlib to create a *sql.DB connection
	db := stdlib.OpenDB(*pgxConfig)
	defer db.Close()

	// Run the migration command
	if err := runMigration(db, direction); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

// runMigration handles the migration based on the specified direction
func runMigration(db *sql.DB, direction string) error {
	// Set up the migration driver
	driver, err := postgresMigrate.WithInstance(db, &postgresMigrate.Config{})
	if err != nil {
		return err
	}

	// Initialize migration with PostgreSQL driver
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Migrations folder path
		"postgres",          // Database name (for logging)
		driver,
	)
	if err != nil {
		return err
	}

	// Display current migration version
	v, dirty, err := m.Version()
	if err == nil {
		log.Printf("Current migration version: %d, dirty: %v", v, dirty)
	} else if err == migrate.ErrNilVersion {
		log.Println("No migrations applied yet.")
	} else {
		return err
	}

	// Run migration based on command
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}
		log.Println("Migration up completed successfully.")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return err
		}
		log.Println("Migration down completed successfully.")
	default:
		log.Println("Unknown direction. Use 'up' or 'down'.")
	}

	return nil
}
