package prodcatalog

import (
	"context"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/domain/dto"
	"github.com/google/uuid"
)

// ProductRepository defines an interface for managing product data, including retrieval, deletion, and restoration operations.
//
//go:generate mockgen -source=service.go -destination mock_product_repository_test.go --package ProductCatalog_test
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

	GetProductLog(ctx context.Context, filter *dto.ProductLogFilter) ([]interface{}, error)
}

// EventBus defines an interface for managing event publishing and subscription.
type EventBus interface {
	// Publish sends an event to all subscribers of the specified topic.
	Publish(ctx context.Context, topic string, event interface{}) error

	// Subscribe registers a handler function to a specific topic.
	// The handler will be called whenever an event is published to that topic.
	Subscribe(ctx context.Context, topic string, handler func(ctx context.Context, event interface{})) error
}

// AuthService returns information about the current command
type AuthService interface {
	GetPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
}

type ProductService struct {
	repo ProductRepository
	bus  EventBus
	auth AuthService
}

func NewProductService(repo ProductRepository, bus EventBus, auth AuthService) (*ProductService, error) {
	return &ProductService{
		repo: repo,
		bus:  bus,
		auth: auth,
	}, nil
}
