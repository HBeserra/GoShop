package ProductCatalog

import (
	"context"
	"github.com/google/uuid"
	"shopper/domain"
	"shopper/domain/dto"
)

// ProductRepository defines an interface for managing product data, including retrieval, deletion, and restoration operations.
//
//go:generate mockgen -source=service.go -destination mock/product_repository_test.go
type ProductRepository interface {

	// Find retrieves a list of products matching the criteria specified in the provided product filter. Returns an error if the operation fails.
	Find(ctx context.Context, filter dto.ProductFilter) ([]*domain.Product, error)

	// GetByID retrieves a product by its unique identifier and returns the product or an error if not found.
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)

	// Update updates the details of an existing product in the repository and returns an error if the operation fails.
	Update(ctx context.Context, product *domain.Product) error

	// Create adds a new product to the repository and returns an error if the operation fails.
	Create(ctx context.Context, product *domain.Product) error

	// Delete removes a product by the provided UUID and returns an error if the operation fails.
	// uses soft-delete.
	Delete(ctx context.Context, id uuid.UUID) error

	// Restore reverts a soft-deleted product by the provided UUID and returns an error if the operation fails.
	Restore(ctx context.Context, id uuid.UUID) error
}

type ProductService struct {
	repo ProductRepository
}
