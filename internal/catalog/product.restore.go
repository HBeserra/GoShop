package catalog

import (
	"context"
	"fmt"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/pkg/observability"
	"github.com/google/uuid"
	"log/slog"
)

func (s *ProductService) RestoreProduct(ctx context.Context, namespace string, productID uuid.UUID) error {

	ctx, span := observability.StartSpan(ctx, "catalog.RestoreProduct")
	defer span.End()

	userID, err := s.auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	perm, err := s.auth.CheckPermissions(ctx, userID, namespace, "product:restore")
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrUnauthorized, err)
	}
	if !perm {
		return domain.ErrUnauthorized
	}

	err = s.repo.Restore(ctx, namespace, productID)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "product restored",
		"product_id", productID,
		"restored_by", userID,
		"namespace", namespace,
	)

	return nil
}
