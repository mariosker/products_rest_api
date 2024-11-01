// Package handlers provides HTTP handlers for managing products.
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mariosker/products_rest_api/internal/models"
	"github.com/mariosker/products_rest_api/internal/repository"
	"github.com/mariosker/products_rest_api/internal/utils"
)

// ProductHandler handles HTTP requests for managing products.
type ProductHandler struct {
	repo repository.ProductRepository
}

// NewProductHandler creates a new ProductHandler with the given repository.
func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.CreateProductPayload true "Product Payload"
// @Success 201 {object} models.CreateProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.CreateProductPayload
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, repoErr := h.repo.CreateProduct(c.Request.Context(), &product)
	if repoErr != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create product")
		return
	}
	response := models.CreateProductResponse{ID: id}
	c.JSON(http.StatusCreated, response)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Retrieve a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID: "+c.Param("id"))
		return
	}

	product, err := h.repo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Product with id: "+strconv.Itoa(id)+" not found")
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts godoc
// @Summary Get a list of products
// @Description Retrieve a list of products with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Product
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid limit")
		return
	}
	if limit <= 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Limit must be greater than 0")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid offset")
		return
	}
	if offset < 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Offset must be greater than or equal to 0")
		return
	}

	products, err := h.repo.GetProducts(c.Request.Context(), limit, offset)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	if products == nil {
		products = []*models.Product{}
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Update an existing product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.UpdateProductPayload true "Product Payload"
// @Success 200 {object} models.Product
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID: "+c.Param("id"))
		return
	}

	var payload models.UpdateProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.repo.UpdateProduct(c.Request.Context(), id, &payload)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with ID %d not found", id) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Product not found")
		} else {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update product with ID: "+strconv.Itoa(id))
		}
		return
	}
	product := models.Product{ID: id, Name: payload.Name, Price: payload.Price}
	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 {} {}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, parseErr := strconv.Atoi(c.Param("id"))
	if parseErr != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID: "+c.Param("id"))
		return
	}

	deleteErr := h.repo.DeleteProduct(c.Request.Context(), id)
	if deleteErr != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete product with id: "+strconv.Itoa(id))
		return
	}

	c.Status(http.StatusNoContent)
}
