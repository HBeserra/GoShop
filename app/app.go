package app

import (
	"context"
	"github.com/HBeserra/GoShop/internal/catalog"
)

type shutdownFn struct {
	Name string
	Func func(ctx context.Context) error
}

type app struct {
	shutdownFn []shutdownFn
	productSvc *catalog.ProductService
}
