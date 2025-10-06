package port

import (
	"context"
	"shop-api-go/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

// AdminService is an interface for interacting with user-related business logic.
type AdminService interface {
	// GetUsersByOffestPagination fetches users using offset pagination.
	GetUsersByOffestPagination(ctx context.Context, token *domain.Token, page, limit int) ([]domain.User, error)
	// GetUsersByTimePagination fetches users using time pagination.
	GetUsersByTimePagination(ctx context.Context, token *domain.Token, after time.Time, limit int) ([]domain.User, error)
	// SearchUserByUsername searches for users with similar to the provided username.
	SearchUserByUsername(ctx context.Context, token *domain.Token, username string, limit int) ([]domain.User, error)
	// SearchUserByEmail searches for users with similar to the provided email.
	SearchUserByEmail(ctx context.Context, token *domain.Token, email string, limit int) ([]domain.User, error) // GetUserById fetches a user by specific.
	// GetUserById fetches a user by specific.
	GetUserById(ctx context.Context, token *domain.Token, id uuid.UUID) (*domain.User, error)
	// UpdateUsername updates the username of specific user by id.
	UpdateUsername(ctx context.Context, token *domain.Token, id uuid.UUID, username string) error
	// UpdateEmail updates the email of specific user by id.
	UpdateEmail(ctx context.Context, token *domain.Token, id uuid.UUID, email string) error
}
