package port

import (
	"context"
	"shop-api-go/internal/core/domain"
	"time"

	"github.com/google/uuid"
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
	// SearchUserByUsername searches for users with similar to the provided username.
	SearchUserByUsername(ctx context.Context, username string, limit int) ([]domain.User, error)
	// SearchUserByEmail searches for users with similar to the provided email.
	SearchUserByEmail(ctx context.Context, email string, limit int) ([]domain.User, error)
	// GetUserById fetches a user by specific.
	GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// UpdateUsername updates the username of specific user by id.
	UpdateUsername(ctx context.Context, id uuid.UUID, username string) error
	// UpdateEmail updates the email of specific user by id.
	UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
	// UpdatePassword updates the password of a specific user by id.
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	// UpdateRole updates the role of a specific user by id.
	UpdateRole(ctx context.Context, id uuid.UUID, role domain.UserRole) error
}

// UserService is an interface for interacting with user-related business logic.
type UserService interface {
	// Register adds a new user.
	Register(ctx context.Context, user *domain.User) error
}
