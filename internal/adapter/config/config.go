package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// getEnv is helper function used to get env variable, if the variable doesn't exist fallback is returned.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvInt is helper function used to get integer env variable,
// if the variable doesn't exist or is invalid integer fallback is returned.
func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

// getEnvDuration is a helper function used to get time.Duration env variable,
// if the variable doesn't exist or is invalid time.Duration fallback is returned.
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := time.ParseDuration(value); err == nil {
			return i
		}
	}
	return fallback
}

type (
	// Environment is an enum for different app environments.
	Environment string
	// Container contains all environment variables.
	Container struct {
		App      *AppConfig
		Database *DBConfig
		JWT      *JWTConfig
	}
	// AppConfig contains all environment variable for the application.
	AppConfig struct {
		Environment Environment
		Port        string
	}

	// DBConfig contains all environment variables for the database.
	DBConfig struct {
		Url                string
		MaxIdleConnections int
		MaxOpenConnections int
	}

	// JWTConfig contains all environment variables for the JWTConfig tokens.
	JWTConfig struct {
		Secret                 []byte
		Issuer                 string
		Audience               string
		RefreshTokenExpireTime time.Duration
		AccessTokenExpireTime  time.Duration
	}
)

const (
	Production  Environment = "production"
	Development Environment = "development"
)

// New creates a new Container instance.
func New() (*Container, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	environment := Environment(getEnv("ENVIRONMENT", string(Development)))
	if environment != Production && environment != Development {
		return nil, fmt.Errorf("unknown environment: %s", environment)
	}

	maxOpenConnections := getEnvInt("DATABASE_MAX_OPEN_CONNECTIONS", 10)
	if maxOpenConnections <= 0 {
		return nil, fmt.Errorf("database max open connections must be > 0: %d", maxOpenConnections)
	}
	maxIdleConnections := getEnvInt("DATABASE_MAX_IDLE_CONNECTIONS", 10)
	if maxIdleConnections <= 0 {
		return nil, fmt.Errorf("database max idle connections must be > 0: %d", maxIdleConnections)
	}

	secret := getEnv("JWT_SECRET", "secret")
	if secret == "" {
		return nil, fmt.Errorf("jwt secret not be empty")
	}
	if secret == "secret" {
		log.Println("WARNING: JWT secret not set, using fallback")
	}

	refreshTokenExpireTime := getEnvDuration("JWT_REFRESH_TOKEN_EXPIRE_TIME", 24*time.Hour)
	if refreshTokenExpireTime <= 0 {
		return nil, fmt.Errorf("jwt refresh token expire time must be > 0: %d", refreshTokenExpireTime)
	}

	accessTokenExpireTime := getEnvDuration("JWT_ACCESS_TOKEN_EXPIRE_TIME", 30*time.Minute)
	if accessTokenExpireTime <= 0 {
		return nil, fmt.Errorf("jwt access token expire time must be > 0: %d", accessTokenExpireTime)
	}

	return &Container{
		App: &AppConfig{
			Environment: environment,
			Port:        getEnv("PORT", "8080"),
		},
		Database: &DBConfig{
			Url:                getEnv("DATABASE_URL", ""),
			MaxOpenConnections: maxOpenConnections,
			MaxIdleConnections: maxIdleConnections,
		},
		JWT: &JWTConfig{
			Secret:                 []byte(secret),
			Issuer:                 getEnv("JWT_ISSUER", "my-app"),
			Audience:               getEnv("JWT_AUDIENCE", "my-app-users"),
			RefreshTokenExpireTime: refreshTokenExpireTime,
			AccessTokenExpireTime:  accessTokenExpireTime,
		},
	}, nil
}
