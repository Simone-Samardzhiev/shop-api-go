package logger

import (
	"shop-api-go/internal/adapter/config"

	"go.uber.org/zap"
)

var logger *zap.Logger

// SetLogger sets Logger based on the passed config.AppConfig.
func SetLogger(appConfig *config.AppConfig) error {
	var err error
	switch appConfig.Environment {
	case config.Production:
		logger, err = zap.NewProduction()
	case config.Development:
		logger, err = zap.NewDevelopment()
	}

	zap.ReplaceGlobals(logger)
	return err
}
