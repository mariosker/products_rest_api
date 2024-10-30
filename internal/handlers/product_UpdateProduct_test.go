package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mariosker/products_rest_api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_UpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	mockRepo := new(MockProductRepository)
	handler := NewProductHandler(mockRepo)

	router.PUT("/products/:id", handler.UpdateProduct)

	t.Run("Product Not Found", func(t *testing.T) {
		var errRepo = errors.New("product with ID 3 not found")
		mockRepo.On("UpdateProduct", mock.Anything, mock.Anything).Return(errRepo).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/products/3", strings.NewReader(`{"name":"Updated Product","price":15.0}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/products/invalid", strings.NewReader(`{"name":"Updated Product","price":15.0}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/products/3", strings.NewReader(`{"name":"Updated Product","price":"invalid"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		product := &models.Product{ID: 3, Name: "Updated Product", Price: 15.0}
		mockRepo.On("UpdateProduct", mock.Anything, product).Return(nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/products/3", strings.NewReader(`{"name":"Updated Product","price":15.0}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Updated Product"`)
	})
}
