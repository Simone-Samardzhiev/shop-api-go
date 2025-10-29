package task

import (
	"context"
	"shop-api-go/internal/core/port"
	"time"

	"go.uber.org/zap"
)

func StartDeleteExpiredTokensTask(ctx context.Context, tokenRepository port.TokenRepository, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				zap.L().Info("deleting expired tokens")
				_ = tokenRepository.DeleteExpiredTokens()
			case <-ctx.Done():
				zap.L().Info("stoping expired token clean up task")
			}
		}
	}()
}
