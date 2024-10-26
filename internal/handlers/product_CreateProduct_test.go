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

func TestProductHandler_CreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	mockRepo := new(MockProductRepository)
	handler := NewProductHandler(mockRepo)

	router.POST("/products", handler.CreateProduct)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("CreateProduct", mock.Anything, &models.CreateProductPayload{Name: "Test Product", Price: 10.0}).Return(1, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(`{"name":"Test Product","price":10.0}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"id":"1"`)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(`{"name":"Test Product", "price":"invalid"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(`{}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Zero Price", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(`{"name":"Test Product","price":0.0}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Negative Price", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(`{"name":"Test Product","price":-10.0}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Creation Error", func(t *testing.T) {
		var errRepo = errors.New("creation error")
		mockRepo.On("CreateProduct", mock.Anything, mock.Anything).Return(-1, errRepo).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products", strings.NewReader(`{"name":"Test Product","price":10.0}`))
		router.ServeHTTP(w, req)

		// Assert that the status code is 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Assert that the response body contains the expected error message
		assert.Contains(t, w.Body.String(), "Failed to create product")
	})

}
