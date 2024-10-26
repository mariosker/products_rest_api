package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mariosker/products_rest_api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_GetProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	mockRepo := new(MockProductRepository)
	handler := NewProductHandler(mockRepo)

	router.GET("/products", handler.GetProducts)

	t.Run("Success", func(t *testing.T) {
		mockProducts := []*models.Product{
			{ID: 1, Name: "Product 1", Price: 10.0},
			{ID: 2, Name: "Product 2", Price: 20.0},
		}
		mockRepo.On("GetProducts", mock.Anything, 10, 0).Return(mockProducts, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products?limit=10&offset=0", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"id":1,"name":"Product 1","price":10},{"id":2,"name":"Product 2","price":20}]`, w.Body.String())
	})
	t.Run("Success with Pagination", func(t *testing.T) {
		mockProducts := []*models.Product{
			{ID: 1, Name: "Product 1", Price: 10.0},
			{ID: 2, Name: "Product 2", Price: 20.0},
			{ID: 3, Name: "Product 3", Price: 30.0},
			{ID: 4, Name: "Product 4", Price: 40.0},
		}
		// First call with limit=2 and offset=0
		mockRepo.On("GetProducts", mock.Anything, 2, 0).Return(mockProducts[:2], nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products?limit=2&offset=0", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"id":1,"name":"Product 1","price":10},{"id":2,"name":"Product 2","price":20}]`, w.Body.String())

		// Second call with limit=2 and offset=2
		mockRepo.On("GetProducts", mock.Anything, 2, 2).Return(mockProducts[2:], nil).Times(1)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/products?limit=2&offset=2", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"id":3,"name":"Product 3","price":30},{"id":4,"name":"Product 4","price":40}]`, w.Body.String())
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockRepo.On("GetProducts", mock.Anything, 10, 0).Return(nil, errors.New("database error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products?limit=10&offset=0", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve products")
	})

	t.Run("Invalid Limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products?limit=invalid&offset=0", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Limit Negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products?limit=-2&offset=0", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Offset Negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products?limit=10&offset=-2", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Offset", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products?limit=10&offset=invalid", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
