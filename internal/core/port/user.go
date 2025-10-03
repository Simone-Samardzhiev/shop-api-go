package port

import (
	"context"
	"shop-api-go/internal/core/domain"
	"time"
)

// UserRepository is an interface for interacting with user-related data.
type UserRepository interface {
	// AddUser inserts a new user into the database.
	AddUser(ctx context.Context, user *domain.User) error
	// GetUserByUsername fetches a user by specific username.
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	// GetUsersByOffestPagination fetches users using offset pagination.
	GetUsersByOffestPagination(ctx context.Context, page, limit int) ([]domain.User, error)
	// GetUsersByTimePagination fetches users using time pagination.
	GetUsersByTimePagination(ctx context.Context, after time.Time, limit int) ([]domain.User, error)
}

// UserService is an interface for interacting with user-related business logic.
type UserService interface {
	// Register adds a new user.
	Register(ctx context.Context, user *domain.User) error
	// GetUsersByOffestPagination fetches users using offset pagination.
	GetUsersByOffestPagination(ctx context.Context, token *domain.Token, page, limit int) ([]domain.User, error)
	// GetUsersByTimePagination fetches users using time pagination.
	GetUsersByTimePagination(ctx context.Context, token *domain.Token, after time.Time, limit int) ([]domain.User, error)
}
