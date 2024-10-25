package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mariosker/products_rest_api/internal/config"
	"github.com/mariosker/products_rest_api/internal/database"
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

	// Set up router and routes
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	serverAddr := cfg.ServerHost + ":" + cfg.ServerPort
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
