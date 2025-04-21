package prodcatalog

import (
	"context"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/domain/dto"
	"github.com/HBeserra/GoShop/pkg/observability"
)

func (s *ProductService) Find(ctx context.Context, namespace string, filter dto.ProductFilter) ([]*domain.Product, error) {

	ctx, span := observability.StartSpan(ctx, "prodcatalog.Find")
	defer span.End()

	userID, err := s.auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	perm, err := s.auth.CheckPermissions(ctx, userID, namespace, "product:read")
	if err != nil {
		return nil, err
	}
	if !perm {
		return nil, domain.ErrUnauthorized
	}

	return s.repo.Find(ctx, filter)
}
