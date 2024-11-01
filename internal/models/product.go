package models

// Product defines the structure for a product
// @Description Product defines the structure for a product
type Product struct {
	ID    int     `json:"id" db:"id"`
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}

// CreateProductPayload defines the payload for creating a product
// @Description CreateProductPayload defines the structure for creating a new product
type CreateProductPayload struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" db:"price" binding:"required,gt=0"`
}

// CreateProductResponse defines the response for creating a product
// @Description CreateProductResponse defines the structure for the response when creating a new product
type CreateProductResponse struct {
	ID int `json:"id"`
}

// UpdateProductPayload defines the payload for updating a product
// @Description UpdateProductPayload defines the structure for updating an existing product
type UpdateProductPayload struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" db:"price" binding:"required,gt=0"`
}
