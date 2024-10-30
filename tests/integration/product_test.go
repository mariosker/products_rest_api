package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/mariosker/products_rest_api/internal/handlers"
	"github.com/mariosker/products_rest_api/internal/repository"
	"github.com/mariosker/products_rest_api/internal/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pgxConn *pgx.Conn

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Set up pgx.Conn for application use
	pgxConn, err = pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect with pgx.Conn: %v", err)
	}

	// Run migrations using a temporary *sql.DB connection
	if err := runMigrations(dsn); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	code := m.Run()
	_ = truncateTables()

	// Clean up resources
	pgxConn.Close(ctx)
	_ = postgresC.Terminate(ctx)
	os.Exit(code)
}

func runMigrations(dsn string) error {
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
		"file://../../migrations",
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

func truncateTables() error {
	_, err := pgxConn.Exec(context.Background(), `
		TRUNCATE TABLE products RESTART IDENTITY CASCADE;
	`)
	return err
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	productRepo := repository.NewPostgresProductRepository(pgxConn)
	productHandler := handlers.NewProductHandler(productRepo)
	routes.SetupRoutes(r, productHandler)
	return r
}

func setupTest(t *testing.T) *gin.Engine {
	require.NoError(t, truncateTables())
	t.Cleanup(func() { _ = truncateTables() })
	return setupRouter()
}

func insertTestProduct(name string, price float64) (int, error) {
	var id int
	query := "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id"
	err := pgxConn.QueryRow(context.Background(), query, name, price).Scan(&id)
	return id, err
}

func TestCreateProduct(t *testing.T) {
	router := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		body := `{"name": "Test Product", "price": 20.0}`
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Negative Price", func(t *testing.T) {
		body := `{"name": "Test Product", "price": -20.0}`
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetProductByID(t *testing.T) {
	router := setupTest(t)

	productID, err := insertTestProduct("Sample Product", 15.0)
	require.NoError(t, err)

	t.Run("Product Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/products/%d", productID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Sample Product"`)
	})
	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("Product Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products/9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Product with id: 9999 not found")
	})
}

func TestUpdateProduct(t *testing.T) {
	router := setupTest(t)

	productID, err := insertTestProduct("Original Product", 10.0)
	require.NoError(t, err)

	t.Run("Update Success", func(t *testing.T) {
		body := `{"name": "Updated Product", "price": 25.0}`
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%d", productID), strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Updated Product"`)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		body := `{"name": "Updated Product", "price": 25.0}`
		req, _ := http.NewRequest("PUT", "/products/invalid", strings.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not existent product", func(t *testing.T) {
		body := `{"name": "Updated Product", "price": 25.0}`
		req, _ := http.NewRequest("PUT", "/products/3", strings.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Only price in JSON", func(t *testing.T) {
		body := `{"price": 25.0}`
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%d", productID), strings.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteProduct(t *testing.T) {
	router := setupTest(t)

	productID, err := insertTestProduct("Product to Delete", 10.0)
	require.NoError(t, err)

	t.Run("Delete Success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/products/%d", productID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Delete Non-Existent Product", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/products/9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestGetProducts(t *testing.T) {
	router := setupTest(t)

	_, _ = insertTestProduct("Product 1", 10.0)
	_, _ = insertTestProduct("Product 2", 20.0)

	t.Run("Get Products with Limit and Offset", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products?limit=1&offset=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Product 2"`)
		assert.NotContains(t, w.Body.String(), `"name":"Product 1"`)
	})

	t.Run("Limit Exceeds Available Products", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products?limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 2, len(response))
	})
}
