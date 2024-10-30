// Package handlers provides HTTP handlers for managing products.
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mariosker/products_rest_api/internal/models"
	"github.com/mariosker/products_rest_api/internal/repository"
)

// ProductHandler handles HTTP requests for managing products.
type ProductHandler struct {
	repo repository.ProductRepository
}

// NewProductHandler creates a new ProductHandler with the given repository.
func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// CreateProduct handles the creation of a new product.
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.CreateProductPayload
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, repoErr := h.repo.CreateProduct(c.Request.Context(), &product)
	if repoErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": strconv.Itoa(id)})
}

// GetProduct handles the HTTP GET request to retrieve a product by its ID.
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID: " + c.Param("id")})
		return
	}

	product, err := h.repo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product with id: " + strconv.Itoa(id) + " not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts handles the HTTP GET request to retrieve a list of products with pagination.
func (h *ProductHandler) GetProducts(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	if limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be greater than 0"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}
	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offset must be greater than or equal to 0"})
		return
	}

	products, err := h.repo.GetProducts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct handles the HTTP PUT request to update an existing product by its ID.
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID: " + c.Param("id")})
		return
	}

	var payload models.UpdateProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.UpdateProduct(c.Request.Context(), id, &payload)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with ID %d not found", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product with ID: " + strconv.Itoa(id)})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "name": payload.Name, "price": payload.Price})
}

// DeleteProduct handles the deletion of a product by its ID.
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, parseErr := strconv.Atoi(c.Param("id"))
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID: " + c.Param("id")})
		return
	}

	deleteErr := h.repo.DeleteProduct(c.Request.Context(), id)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product with id: " + strconv.Itoa(id)})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
