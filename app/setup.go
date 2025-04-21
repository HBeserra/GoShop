package app

import (
	"context"
	"github.com/HBeserra/GoShop/internal/prodcatalog"
	"log/slog"
)

func New(ctx context.Context) {

	var a = new(app)

	/*
	 *	Get the config
	 */

	/*
	 *	Set up the repositories
	 */

	/*
	 *	Set up the Services
	 */

	// Set up the Product Catalog Service
	prodSvc, err := prodcatalog.NewProductService(nil, nil, nil)
	a.ifErrShutdown(ctx, err)
	a.productSvc = prodSvc

	/*
	 *	Start the controllers
	 */
}

func (a *app) Shutdown(ctx context.Context) error {

	for i := len(a.shutdownFn) - 1; i >= 0; i-- {
		fn := a.shutdownFn[i]
		if err := ctx.Err(); err != nil {
			slog.ErrorContext(ctx, "shutdown error", "error", err)
			return ctx.Err()
		}

		if err := fn.Func(ctx); err != nil {
			slog.ErrorContext(ctx, "shutdown error", "error", err, "fn", fn.Name)
			return err
		}
	}

	return nil
}

func (a *app) addShutdownFn(name string, fn func(ctx context.Context) error) {
	a.shutdownFn = append(a.shutdownFn, shutdownFn{
		Name: name,
		Func: fn,
	})
}

func (a *app) ifErrShutdown(ctx context.Context, err error) {
	if err != nil {
		a.Shutdown(ctx)
	}
}
