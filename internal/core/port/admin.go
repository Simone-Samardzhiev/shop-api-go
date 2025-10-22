package port

import (
	"context"
	"shop-api-go/internal/core/domain"
)

// AdminService is an interface for interacting with user-related business logic.
type AdminService interface {
	// GetUsers fetches users by passed filters.
	GetUsers(ctx context.Context, token *domain.Token, get *domain.GetUsers) (*domain.UsersResult, error)
	// UpdateUser updates a specific user field.
	UpdateUser(ctx context.Context, token *domain.Token, update *domain.UserUpdate) error
}
