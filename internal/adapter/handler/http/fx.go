package http

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"HTTP",
	fx.Provide(NewUserHandler),
	fx.Provide(NewAuthHandler),
	fx.Provide(NewAdminHandler),
	fx.Provide(NewRouter),
	fx.Invoke(func(lc fx.Lifecycle, router *Router) {
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				var err error
				go func() {
					err = router.Start()
				}()
				return err
			},
			OnStop: func(ctx context.Context) error {
				return router.Shutdown(ctx)
			},
		})
	}),
)
