package dto

import (
	"github.com/google/uuid"
	"shopper/domain"
	"time"
)

type ProductFilter struct {
	IDs []uuid.UUID `json:"ids"`
}

type ProductLogFilter struct {
	Namespace string
	ProductID []uuid.UUID           `json:"product_id"`
	Events    []domain.ProductEvent `json:"events"`
	UserID    []uuid.UUID           `json:"user_id"`

	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Offset int       `json:"offset"`
	Limit  int       `json:"limit"`
}
