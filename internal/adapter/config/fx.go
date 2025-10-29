package config

import "go.uber.org/fx"

var Module = fx.Module(
	"config",
	fx.Provide(New),
	fx.Provide(func(config *Container) *AppConfig {
		return config.App
	}),
	fx.Provide(func(config *Container) *DBConfig {
		return config.Database
	}),
	fx.Provide(func(config *Container) *JWTConfig {
		return config.JWT
	}),
)
