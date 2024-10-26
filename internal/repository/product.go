package repository

import (
	"context"

	"github.com/mariosker/products_rest_api/internal/database"
	"github.com/mariosker/products_rest_api/internal/models"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.CreateProductPayload) (int, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*models.Product, error)
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

func (r *PostgresProductRepository) GetProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, price FROM products LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
