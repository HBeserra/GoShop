package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	// ID is the unique identifier for the ProductVariant and serves as the primary key.
	ID uuid.UUID `json:"id" gorm:"primaryKey;index:idx_user"`
	// Namespace specifies the logical grouping in the multi-tenant system
	Namespace string `json:"namespace" gorm:"index:idx_user"`
	// CreatedAt indicates the timestamp when the entity was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt indicates the timestamp of the most recent update to the entity.
	UpdatedAt time.Time `json:"updated_at"`

	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex;index:idx_user"`
}
