package catalog

import (
	"context"
	"fmt"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/pkg/observability"
	"github.com/google/uuid"
	"log/slog"
)

func (s *ProductService) DeleteProduct(ctx context.Context, namespace string, productID uuid.UUID) error {

	ctx, span := observability.StartSpan(ctx, "catalog.DeleteProduct")
	defer span.End()

	userID, err := s.auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	perm, err := s.auth.CheckPermissions(ctx, userID, namespace, "product:delete")
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrUnauthorized, err)
	}
	if !perm {
		return domain.ErrUnauthorized
	}

	err = s.repo.Delete(ctx, namespace, productID)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "product deleted",
		"product_id", productID,
		"deleted_by", userID,
	)

	return nil
}
