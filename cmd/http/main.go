package main

import (
	"shop-api-go/internal/adapter/auth"
	"shop-api-go/internal/adapter/config"
	"shop-api-go/internal/adapter/handler/http"
	"shop-api-go/internal/adapter/logger"
	"shop-api-go/internal/adapter/storage/postgres"
	"shop-api-go/internal/core/service"
	"shop-api-go/internal/core/task"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

// @title Shop API
// @version 1.0
// @description     This is the Shop API server for e-shop.
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http
func main() {
	fx.New(
		config.Module,
		logger.Module,
		postgres.Module,
		auth.Module,
		service.Module,
		task.Module,
		http.Module,
	).Run()
}
