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

func TestProductHandler_GetProductByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	mockRepo := new(MockProductRepository)
	handler := NewProductHandler(mockRepo)

	router.GET("/products/:id", handler.GetProduct)

	t.Run("Success", func(t *testing.T) {
		mockProduct := &models.Product{ID: 1, Name: "Test Product", Price: 10.0}
		mockRepo.On("GetProductByID", mock.Anything, 1).Return(mockProduct, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Test Product"`)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		var errRepo = errors.New("product not found")
		mockRepo.On("GetProductByID", mock.Anything, 3).Return(nil, errRepo).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products/3", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Product with id: 3 not found")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products/invalid_id", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
