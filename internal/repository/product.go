package repository

import (
	"context"
	"fmt"

	"github.com/mariosker/products_rest_api/internal/database"
	"github.com/mariosker/products_rest_api/internal/models"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.CreateProductPayload) (int, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, id int, payload *models.UpdateProductPayload) error
	DeleteProduct(ctx context.Context, id int) error
}

type PostgresProductRepository struct {
	dbConnection database.DBConnection
}

func NewPostgresProductRepository(dbConnection database.DBConnection) *PostgresProductRepository {
	return &PostgresProductRepository{dbConnection: dbConnection}
}

// CreateProduct inserts a new product into the database and returns the new product's ID.
// Parameters:
// - ctx: The context for managing request-scoped values, cancelation, and deadlines.
// - product: The payload containing the product details to be created.
func (r *PostgresProductRepository) CreateProduct(ctx context.Context, product *models.CreateProductPayload) (int, error) {
	var id int

	err := r.dbConnection.QueryRow(ctx, "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id", product.Name, product.Price).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// GetProductByID retrieves a product from the database by its ID.
// Parameters:
// - ctx: context for managing request deadlines and cancellation signals.
// - id: the ID of the product to be retrieved.
func (r *PostgresProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	err := r.dbConnection.QueryRow(ctx, "SELECT id, name, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProducts retrieves a list of products from the database with pagination support.
// Parameters:
// - ctx: context for managing request deadlines and cancellation signals.
// - limit: the maximum number of products to return.
// - offset: the number of products to skip before starting to return products.
func (r *PostgresProductRepository) GetProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	rows, err := r.dbConnection.Query(ctx, "SELECT id, name, price FROM products LIMIT $1 OFFSET $2", limit, offset)
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

// UpdateProduct updates an existing product in the database.
// Parameters:
// - ctx: context for managing request deadlines and cancellation signals.
// - id: the ID of the product to be updated.
// - payload: the product data to be updated.
func (r *PostgresProductRepository) UpdateProduct(ctx context.Context, id int, payload *models.UpdateProductPayload) error {
	result, err := r.dbConnection.Exec(ctx, "UPDATE products SET name=$1, price=$2 WHERE id=$3", payload.Name, payload.Price, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}

// DeleteProduct deletes a product from the database by its ID.
// Parameters:
// - ctx: context for managing request deadlines and cancellation signals.
// - id: the ID of the product to be deleted.
func (r *PostgresProductRepository) DeleteProduct(ctx context.Context, id int) error {
	_, err := r.dbConnection.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
