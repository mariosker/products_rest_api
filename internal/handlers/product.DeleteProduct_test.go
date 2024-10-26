package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_DeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	mockRepo := new(MockProductRepository)
	handler := NewProductHandler(mockRepo)

	router.DELETE("/products/:id", handler.DeleteProduct)

	t.Run("Product Not Found", func(t *testing.T) {
		var errRepo = errors.New("product not found")
		mockRepo.On("DeleteProduct", mock.Anything, 3).Return(errRepo).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/products/3", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to delete product with id: 3")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/products/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteProduct", mock.Anything, 3).Return(nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/products/3", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
