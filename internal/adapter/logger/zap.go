package logger

import (
	"fmt"
	"shop-api-go/internal/adapter/config"

	"go.uber.org/zap"
)

var logger *zap.Logger

// SetLogger sets Logger based on the passed config.AppConfig.
func SetLogger(config *config.AppConfig) error {
	var err error
	switch config.Environment {
	case "development":
		logger, err = zap.NewDevelopment()
	case "production":
		logger, err = zap.NewProduction()
	default:
		err = fmt.Errorf("unknown logger environment: %s", config.Environment)
	}

	zap.ReplaceGlobals(logger)
	return err
}
