package util

import (
	"shop-api-go/internal/core/port"
	"time"

	"go.uber.org/zap"
)

func StartDeleteExpiredTokensTask(tokenRepository port.TokenRepository, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			zap.L().Info("Deleting expired tokens")
			err := tokenRepository.DeleteExpiredTokens()
			if err != nil {
				zap.L().Error("Error deleting expired tokens", zap.Error(err))
			}
		}
	}()
}
