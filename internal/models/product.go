package models

type Product struct {
	ID    int     `json:"id" db:"id"`
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}

type CreateProductPayload struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" db:"price" binding:"required,gt=0"`
}

type UpdateProductPayload struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" db:"price" binding:"required,gt=0"`
}
