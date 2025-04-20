package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"shopper/pkg/currency"
)

type Product struct {
	BaseDomain
	Name  string       `json:"name"`
	Price currency.BRL `json:"price"`
	Stock int64        `json:"stock"`
	// Medias represents a collection of associated media objects for a product, stored with many-to-many relationship mapping.
	Medias []Media `json:"medias" gorm:"embedded;many2many:product_medias;"`
	// Variants represents a collection of associated variant objects for a product, such as size or color options.
	Variants []ProductVariant
}

// ProductVariant represents a distinct variation of a product, including attributes such as price, stock, and name.
// It is associated with a specific product via the ProductID field.
// This type extends gorm.Model, providing standard fields like ID, CreatedAt, and UpdatedAt.
type ProductVariant struct {
	gorm.Model
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`

	// ProductID represents the unique identifier of the product associated with this variant.
	ProductID uuid.UUID `json:"product_id"`
	// Name specifies the name of the product variant.
	Name string `json:"name"`
	// Price represents the cost of the product variant as a floating-point number.
	Price currency.BRL `json:"price"`
	// Stock indicates the quantity of this product variant available in inventory.
	Stock int64 `json:"stock"`
}

type BaseDomain struct {
	gorm.Model
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
