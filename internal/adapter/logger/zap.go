package logger

import (
	"fmt"
	"shop-api-go/internal/adapter/config"

	"go.uber.org/zap"
)

// New sets Logger based on the passed config.AppConfig.
func New(appConfig *config.AppConfig) (*zap.Logger, error) {
	switch appConfig.Environment {
	case config.Production:
		cnf := zap.NewProductionConfig()
		cnf.DisableStacktrace = true
		return cnf.Build()
	case config.Development:
		return zap.NewDevelopment()
	default:
		return nil, fmt.Errorf("unknown environment: %s", appConfig.Environment)
	}
}
