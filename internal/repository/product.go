package repository

import (
	"context"

	"github.com/mariosker/products_rest_api/internal/database"
	"github.com/mariosker/products_rest_api/internal/models"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.CreateProductPayload) (int, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
}

type PostgresProductRepository struct {
	db database.DBConnection
}

func NewPostgresProductRepository(db database.DBConnection) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) CreateProduct(ctx context.Context, product *models.CreateProductPayload) (int, error) {
	var id int

	err := r.db.QueryRow(ctx, "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id", product.Name, product.Price).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *PostgresProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	err := r.db.QueryRow(ctx, "SELECT id, name, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}