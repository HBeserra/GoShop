package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"shopper/pkg/currency"
	"time"
)

// Product represents a product entity with details such as ID, name, price, stock, and associated media and variants.
type Product struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"primaryKey;index:idx_product"` // ID is the primary key for the Product entity and is uniquely indexed.
	Namespace string    `json:"namespace" gorm:"index:idx_product"`     // Namespace specifies the logical grouping in the multi-tenant system
	CreatedAt time.Time `json:"created_at"`                             // CreatedAt indicates the timestamp when the entity was created, stored as a string in the JSON response.
	UpdatedAt time.Time `json:"updated_at"`                             // UpdatedAt indicates the last time the entity was modified

	Title  string        `json:"title"`
	Price  currency.BRL  `json:"price"`
	Stock  int64         `json:"stock"`
	Status ProductStatus `json:"status"`
	// SKU is the stock-keeping unit, a unique identifier for inventory tracking.
	SKU string `json:"sku" gorm:"index:idx_product"`
	// Medias represents a collection of associated media objects for a product, stored with many-to-many relationship mapping.
	Medias []Media `json:"medias" gorm:"many2many:product_medias;"`
	// Variants represents a collection of associated variant objects for a product, such as size or color options.
	Variants []ProductVariant
}

// ProductVariant represents a distinct variation of a product, including attributes such as price, stock, and name.
// It is associated with a specific product via the ProductID field.
// This type extends gorm.Model, providing standard fields like ID, CreatedAt, and UpdatedAt.
type ProductVariant struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"primaryKey;index:idx_product_variant"` // ID is the unique identifier for the ProductVariant and serves as the primary key.
	Namespace string    `json:"namespace" gorm:"index:idx_product_variant"`     // Namespace specifies the logical grouping in the multi-tenant system
	CreatedAt time.Time `json:"created_at"`                                     // CreatedAt indicates the timestamp when the entity was created
	UpdatedAt time.Time `json:"updated_at"`                                     // UpdatedAt indicates the timestamp of the last update to this entity.

	// ProductID represents the unique identifier of the product associated with this variant.
	ProductID uuid.UUID `json:"product_id" gorm:"index:idx_product_variant"`
	// Title specifies the name of the product variant.
	Title string `json:"title"`
	// Price represents the cost of the product variant as a floating-point number.
	Price currency.BRL `json:"price"`
	// Stock indicates the quantity of this product variant available in inventory.
	Stock int64 `json:"stock"`
	// Medias represents a collection of associated media for the product variant, using a many-to-many relationship.
	Medias []Media `json:"medias" gorm:"many2many:product_variant_medias;"`

	ShortDesc string `json:"short_desc"`
	HtmlDesc  string `json:"html_desc"`
	TextDesc  string `json:"text_desc"`
}

type ProductLogEvent struct {
	gorm.Model
	Namespace string       `json:"namespace" gorm:"index:idx_product_log_event"`
	ProductID uuid.UUID    `json:"product_id" gorm:"primaryKey,idx_product_log_event"`
	Timestamp time.Time    `json:"timestamp" gorm:"primaryKey,idx_product_log_event"`
	Event     ProductEvent `json:"event" gorm:"index:idx_product_log_event"`
	Data      interface{}  `json:"data"`
	UserID    uuid.UUID    `json:"user_id" gorm:"index:idx_product_log_event"`
}

type ProductStatus string

const (
	ProductStatusDraft      ProductStatus = "draft"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
	ProductStatusAvailable  ProductStatus = "available"
)

type ProductEvent string

const (
	ProductCreated     ProductEvent = "product.created"
	ProductUpdated     ProductEvent = "product.updated"
	ProductDeleted     ProductEvent = "product.deleted"
	ProductRestored    ProductEvent = "product.restored"
	ProductStockUpdate ProductEvent = "product.stock.update"
)
