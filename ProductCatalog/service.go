package ProductCatalog

import (
	"context"
	"github.com/google/uuid"
	"shopper/domain"
)

// ProductRepository defines an interface for managing product data, including retrieval, deletion, and restoration operations.
type ProductRepository interface {

	// GetByID retrieves a product by its unique identifier and returns the product or an error if not found.
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)

	// Delete removes a product by the provided UUID and returns an error if the operation fails.
	// uses soft-delete.
	Delete(ctx context.Context, id uuid.UUID) error

	// Restore reverts a soft-deleted product by the provided UUID and returns an error if the operation fails.
	Restore(ctx context.Context, id uuid.UUID) error
}
