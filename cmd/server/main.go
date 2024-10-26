package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mariosker/products_rest_api/internal/config"
	"github.com/mariosker/products_rest_api/internal/database"
	"github.com/mariosker/products_rest_api/internal/handlers"
	"github.com/mariosker/products_rest_api/internal/repository"
	"github.com/mariosker/products_rest_api/internal/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize the database connection
	err = database.InitDB(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	if err := runMigrations(cfg.DBURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Create repository and handlers
	productRepo := repository.NewPostgresProductRepository(database.GetDB())
	productHandler := handlers.NewProductHandler(productRepo)

	// Set up router and routes
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	routes.SetupRoutes(r, productHandler)

	serverAddr := cfg.ServerHost + ":" + cfg.ServerPort
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to run server:", err)
	}
	log.Printf("Running server at: %s", serverAddr)
}

func runMigrations(dsn string) error {
	// Create a temporary *sql.DB connection for migrations
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect for migrations: %w", err)
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("migration setup error: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up error: %w", err)
	}
	return nil
}
