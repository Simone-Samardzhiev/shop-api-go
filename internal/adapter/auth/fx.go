package auth

import (
	"shop-api-go/internal/adapter/auth/bcrypt"
	"shop-api-go/internal/adapter/auth/jwt"
	"shop-api-go/internal/core/port"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"Auth",
	fx.Provide(
		fx.Annotate(
			bcrypt.NewPasswordHasher,
			fx.As(new(port.PasswordHasher)),
		),
	),
	fx.Provide(
		fx.Annotate(
			jwt.NewTokenGenerator,
			fx.As(new(port.TokenGenerator)),
		),
	),
)
