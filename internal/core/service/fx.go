package service

import (
	"shop-api-go/internal/core/port"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"Service",
	fx.Provide(
		fx.Annotate(
			NewUserService,
			fx.As(new(port.UserService)),
		),
	),
	fx.Provide(
		fx.Annotate(
			NewAuthService,
			fx.As(new(port.AuthService)),
		),
	),
	fx.Provide(
		fx.Annotate(
			NewAdminService,
			fx.As(new(port.AdminService)),
		),
	),
)
