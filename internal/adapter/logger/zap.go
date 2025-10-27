package logger

import (
	"fmt"
	"shop-api-go/internal/adapter/config"

	"go.uber.org/zap"
)

var logger *zap.Logger

// New sets Logger based on the passed config.AppConfig.
func New(appConfig *config.AppConfig) (*zap.Logger, error) {
	switch appConfig.Environment {
	case config.Production:
		return zap.NewProduction()
	case config.Development:
		return zap.NewDevelopment()
	default:
		return nil, fmt.Errorf("unknown environment: %s", appConfig.Environment)
	}
}
