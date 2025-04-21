package catalog

import (
	"context"
	"fmt"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/pkg/observability"
	"log/slog"
	"time"
)

func (s *ProductService) UpdateProduct(
	ctx context.Context,
	namespace string,
	product *domain.Product,
) error {

	ctx, span := observability.StartSpan(ctx, "UpdateProduct")
	defer span.End()

	userID, err := s.auth.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrUnauthorized, err)
	}

	perm, err := s.auth.CheckPermissions(ctx, userID, namespace, "product:update")
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrUnauthorized, err)
	}
	if !perm {
		return domain.ErrUnauthorized
	}

	err = s.Validate(ctx, namespace, product)
	if err != nil {
		return err
	}

	product.UpdatedAt = time.Now()
	err = s.repo.Update(ctx, namespace, product)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update product",
			"product_id", product.ID,
			"updated_by", userID,
			"error", err,
		)
		return err
	}

	slog.InfoContext(ctx, "product updated",
		"product_id", product.ID,
		"updated_by", userID,
	)

	return nil

}
