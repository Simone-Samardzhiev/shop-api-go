package port

import (
	"context"
	"shop-api-go/internal/core/domain"
)

// UserRepository is an interface for interacting with user-related data.
type UserRepository interface {
	// CreateUser inserts a new user into the database.
	CreateUser(ctx context.Context, user *domain.User) error
}

// UserService is an interface for interacting with user-related business logic.
type UserService interface {
	Register(ctx context.Context, user *domain.User) error
}
