package prodcatalog

import (
	"context"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/pkg/currency"
	"github.com/HBeserra/GoShop/pkg/observability"
	"slices"
)

func (s *ProductService) Validate(ctx context.Context, product *domain.Product) error {

	ctx, span := observability.StartSpan(ctx, "prodcatalog.Validate")
	defer span.End()

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
	return nil
}
