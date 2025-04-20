package dto

import "github.com/google/uuid"

type ProductFilter struct {
	IDs []uuid.UUID `json:"ids"`
}
