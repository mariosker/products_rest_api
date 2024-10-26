package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mariosker/products_rest_api/internal/handlers"
)

func SetupRoutes(r *gin.Engine, productHandler *handlers.ProductHandler) {
	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products/:id", productHandler.GetProduct)
	r.GET("/products", productHandler.GetProducts)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	// TODO:
	// r.DELETE("/products/:id", productHandler.DeleteProduct)
}
