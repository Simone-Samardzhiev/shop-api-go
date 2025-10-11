package tasks

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
				zap.L().Info("Deleting expired tokens")
				err := tokenRepository.DeleteExpiredTokens()
				if err != nil {
					zap.L().Error("Error deleting expired tokens", zap.Error(err))
				}
			case <-ctx.Done():
				zap.L().Info("Stoping expired token clean up task")
			}
		}
	}()
}
