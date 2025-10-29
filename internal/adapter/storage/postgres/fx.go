package postgres

import (
	"shop-api-go/internal/adapter/storage/postgres/repository"
	"shop-api-go/internal/core/port"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"Postgres",
	fx.Provide(New),
	fx.Provide(
		fx.Annotate(
			repository.NewUserRepository,
			fx.As(new(port.UserRepository)),
		),
	),
	fx.Provide(
		fx.Annotate(
			repository.NewTokenRepository,
			fx.As(new(port.TokenRepository)),
		),
	),
)
