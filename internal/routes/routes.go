package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mariosker/products_rest_api/internal/handlers"
)

func SetupRoutes(r *gin.Engine, productHandler *handlers.ProductHandler) {
	r.POST("/products", productHandler.CreateProduct)
	// TODO:
	// r.GET("/products/:id", productHandler.GetProduct)
	// r.PUT("/products/:id", productHandler.UpdateProduct)
	// r.DELETE("/products/:id", productHandler.DeleteProduct)
	// r.GET("/products", productHandler.GetProducts)
}
