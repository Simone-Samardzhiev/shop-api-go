package main

import (
	"context"
	"log"
	"shop-api-go/internal/adapter/auth/bcrypt"
	"shop-api-go/internal/adapter/auth/jwt"
	"shop-api-go/internal/adapter/config"
	"shop-api-go/internal/adapter/handler/http"
	"shop-api-go/internal/adapter/logger"
	"shop-api-go/internal/adapter/storage/postgres"
	"shop-api-go/internal/adapter/storage/postgres/repository"
	"shop-api-go/internal/core/service"
	"shop-api-go/internal/core/util"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	container, err := config.New()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = logger.SetLogger(container.App)
	if err != nil {
		log.Fatalf("Error setting logger: %v", err)
	}

	db, err := postgres.New(container.Database)
	if err != nil {
		zap.L().Fatal("Error connecting to database", zap.Error(err))
	}

	zap.L().Info("Environment variables",
		zap.Dict("appConfig",
			zap.String("environment", string(container.App.Environment)),
			zap.String("port", container.App.Port),
		),
		zap.Dict("dbConfig",
			zap.Int("maxIdleConnections", container.Database.MaxIdleConnections),
			zap.Int("maxOpenConnections", container.Database.MaxOpenConnections),
		),
		zap.Dict("jwtConfig",
			zap.String("issuer", container.JWT.Issuer),
			zap.String("audience", container.JWT.Audience),
			zap.Duration("refreshTokenExpireTime", container.JWT.RefreshTokenExpireTime),
			zap.Duration("accessTokenExpireTime", container.JWT.AccessTokenExpireTime),
		),
	)

	passwordHasher := &bcrypt.PasswordHasher{}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, passwordHasher)
	userHandler := http.NewUserHandler(userService)

	adminService := service.NewAdminService(userRepository, passwordHasher)
	adminHandler := http.NewAdminHandler(adminService)

	tokenRepository := repository.NewTokenRepository(db)
	jwtTokenGenerator := jwt.NewTokenGenerator(container.JWT)
	authService := service.NewAuthService(jwtTokenGenerator, passwordHasher, tokenRepository, userRepository)
	authHandler := http.NewAuthHandler(authService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	util.StartDeleteExpiredTokensTask(ctx, tokenRepository, time.Hour)

	router := http.NewRouter(container.App, jwtTokenGenerator, userHandler, adminHandler, authHandler)
	err = router.Start()
	if err != nil {
		zap.L().Error("Error starting http server", zap.Error(err))
	}
}
