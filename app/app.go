package app

import (
	"context"
	"github.com/HBeserra/GoShop/internal/prodcatalog"
)

type shutdownFn struct {
	Name string
	Func func(ctx context.Context) error
}

type app struct {
	shutdownFn []shutdownFn
	productSvc *prodcatalog.ProductService
}
