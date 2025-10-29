package logger

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"logger",
	fx.Provide(NewZapLogger),
	fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
		lc.Append(
			fx.Hook{
				OnStart: func(context.Context) error {
					zap.ReplaceGlobals(logger)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return logger.Sync()
				},
			})
	}),
)
