package events

import (
	"github.com/google/uuid"
	"time"
)

type ProductCreated struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedOn time.Time `json:"created_on"`
	CreatedBy uuid.UUID `json:"created_by"`
}

type ProductUpdated struct {
	ID        uuid.UUID `json:"id"`
	UpdatedOn time.Time `json:"updated_on"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}
