package task

import (
	"context"
	"shop-api-go/internal/core/port"
	"time"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"Task",
	fx.Invoke(func(lc fx.Lifecycle, repository port.TokenRepository) {
		bgCtx, cancel := context.WithCancel(context.Background())
		StartDeleteExpiredTokensTask(bgCtx, repository, time.Hour)

		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				cancel()
				return nil
			},
		})
	}),
)
