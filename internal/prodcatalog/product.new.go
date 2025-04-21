package prodcatalog

import (
	"context"
	"fmt"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/domain/events"
	"github.com/HBeserra/GoShop/pkg/currency"
	"github.com/HBeserra/GoShop/pkg/observability"
	"github.com/google/uuid"
	"log/slog"
	"slices"
	"time"
)

func (s *ProductService) CreateProduct(
	ctx context.Context,
	namespace string,
	product *domain.Product,
) error {

	ctx, span := observability.StartSpan(ctx, "ProductService.CreateProduct")
	defer span.End()

	userID, err := s.auth.GetUserID(ctx)
	if err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.ErrUnauthorized
	}

	perm, err := s.auth.CheckPermissions(ctx, userID, namespace, "product:create")
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrUnauthorized, err)
	}
	if !perm {
		return domain.ErrUnauthorized
	}

	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if len(product.Title) < 10 || len(product.Title) > 100 {
		return domain.ErrInvalidProductTitle
	}

	if product.Price.LessOrEqual(currency.NewFromFloat(1.00)) {
		return domain.ErrInvalidProductPrice
	}

	var minPrice, maxPrice float64
	for _, variant := range product.Variants {
		if variant.Price.LessOrEqual(currency.NewFromFloat(1.00)) {
			return domain.ErrInvalidProductPrice
		}

		if variant.Price.Float64() < minPrice {
			minPrice = variant.Price.Float64()
		}
		if variant.Price.Float64() > maxPrice {
			maxPrice = variant.Price.Float64()
		}
	}

	if maxPrice/minPrice > 5 {
		return domain.ErrInvalidProductPrice
	}

	if product.Stock < 0 {
		product.Stock = 0
	}

	if !slices.Contains([]domain.ProductStatus{
		domain.ProductStatusDraft,
		domain.ProductStatusAvailable,
		domain.ProductStatusOutOfStock,
	}, product.Status) {
		return domain.ErrInvalidProductStatus
	}

	// Validate the medias
	for _, mediaID := range product.Medias {
		_, err := s.media.GetByID(ctx, mediaID)
		if err != nil {
			return domain.ErrInvalidMedia
		}
	}

	for _, variant := range product.Variants {
		for _, mediaID := range variant.Medias {
			_, err := s.media.GetByID(ctx, mediaID)
			if err != nil {
				return domain.ErrInvalidMedia
			}
		}
	}

	err = s.repo.Create(ctx, product)
	if err != nil {
		span.RecordError(err)
		return domain.ErrFailedToCreateProduct
	}

	err = s.bus.Publish(ctx, "product:created", events.ProductCreated{
		ID:        product.ID,
		Title:     product.Title,
		CreatedOn: product.CreatedAt,
		CreatedBy: userID,
	})
	if err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "failed to publish product:created event",
			"product_id", product.ID,
			"title", product.Title,
			"created_on", product.CreatedAt,
			"created_by", userID,
			"error", err)
		return nil
	}

	slog.InfoContext(ctx, "product created",
		"product_id", product.ID,
		"title", product.Title,
		"created_on", product.CreatedAt,
		"created_by", userID,
	)
	return nil
}
