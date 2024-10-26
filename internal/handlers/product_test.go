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
