package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// getEnv is helper function used to get env variable, if the variable doesn't exist fallback is returned.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvInt is helper function used to get env variable,
// if the variable doesn't exist or is invalid integer fallback is returned.
func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

type (
	// Container contains all environment variables.
	Container struct {
		Database *Database
	}
	// Database contains all environment variables for the database.
	Database struct {
		Url                string
		MaxIdleConnection  int
		MaxOpenConnections int
	}
)

// New creates a new container instance.
func New() (*Container, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Container{
		Database: &Database{
			Url:                getEnv("DATABASE_URL", ""),
			MaxOpenConnections: getEnvInt("DATABASE_MAX_OPEN_CONNECTIONS", 10),
			MaxIdleConnection:  getEnvInt("DATABASE_MAX_IDLE_CONNECTIONS", 10),
		},
	}, nil
}
