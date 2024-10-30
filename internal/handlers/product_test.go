package handlers

import (
	"context"

	"github.com/mariosker/products_rest_api/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of ProductRepository.
// It is used to simulate the behavior of the actual ProductRepository in tests,
// allowing for controlled responses and verification of interactions without
// needing a real database or external dependencies.
type MockProductRepository struct {
	mock.Mock
}

// CreateProduct mocks the creation of a new product in the repository.
// It takes a context and a CreateProductPayload, and returns the ID of the created product and an error if any.
func (m *MockProductRepository) CreateProduct(ctx context.Context, product *models.CreateProductPayload) (int, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(int), args.Error(1)
}

// GetProductByID retrieves a product by its ID from the mock repository.
// It takes a context and an integer ID as parameters and returns a Product and an error.
func (m *MockProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if product, ok := args.Get(0).(*models.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetProducts mocks the retrieval of a list of products from the repository.
// It takes a context, a limit for the number of products to retrieve, and an offset for pagination.
// It returns a slice of Product pointers and an error if any.
func (m *MockProductRepository) GetProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	args := m.Called(ctx, limit, offset)
	if products, ok := args.Get(0).([]*models.Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateProduct mocks the update of a product in the repository.
// It takes a context and a Product, and returns an error if any.
func (m *MockProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

// DeleteProduct mocks the deletion of a product in the repository.
// It takes a context and an ID of the product to be deleted, and returns an error if any.
func (m *MockProductRepository) DeleteProduct(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
