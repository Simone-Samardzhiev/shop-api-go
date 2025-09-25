package main

import (
	"log"
	"shop-api-go/internal/adapter/config"
	"shop-api-go/internal/adapter/handler/http"
	"shop-api-go/internal/adapter/logger"
	"shop-api-go/internal/adapter/storage/postgres"
	"shop-api-go/internal/adapter/storage/postgres/repository"
	"shop-api-go/internal/core/service"

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

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := http.NewUserHandler(userService)

	router := http.NewRouter(userHandler)
	err = router.Start(container.App.Port)
	if err != nil {
		zap.L().Error("Error starting http server", zap.Error(err))
	}
}
