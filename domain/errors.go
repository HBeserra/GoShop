package domain

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// Product related errors
var (
	ErrProductNotFound       = errors.New("product not found")
	ErrInvalidProductPrice   = errors.New("invalid product price")
	ErrInvalidProductTitle   = errors.New("invalid product title")
	ErrInvalidProductStatus  = errors.New("invalid product status")
	ErrFailedToCreateProduct = errors.New("failed to create product")
	ErrInvalidMedia          = errors.New("invalid media")
)
