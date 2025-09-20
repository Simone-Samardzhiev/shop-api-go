package main

import (
	"log"
	"shop-api-go/internal/adapter/config"
	"shop-api-go/internal/adapter/handler/http"
	"shop-api-go/internal/adapter/storage"
	"shop-api-go/internal/adapter/storage/repository"
	"shop-api-go/internal/core/service"

	_ "github.com/lib/pq"
)

func main() {
	container, err := config.New()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := storage.New(container.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := http.NewUserHandler(userService)

	router := http.NewRouter(userHandler)
	err = router.Start(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
