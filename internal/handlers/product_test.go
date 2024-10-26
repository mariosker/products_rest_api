package handlers

import (
	"context"

	"github.com/mariosker/products_rest_api/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product *models.CreateProductPayload) (int, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if product, ok := args.Get(0).(*models.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}
