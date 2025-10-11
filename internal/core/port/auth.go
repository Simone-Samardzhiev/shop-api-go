package port

import (
	"context"
	"shop-api-go/internal/core/domain"

	"github.com/google/uuid"
)

// TokenGenerator is an interface for signing and decoding tokens.
type TokenGenerator interface {
	// SignToken returns a signed string token.
	SignToken(token *domain.Token) (string, error)
	// ParseToken takes a signed token and returns *domain.Token
	ParseToken(token string) (*domain.Token, error)
}

// PasswordHasher is an interface for hashing and validating passwords.
type PasswordHasher interface {
	// Hash returns hashed version of the password.
	Hash(password string) (string, error)
	// Compare validates that the password and the hash match.
	Compare(password, hash string) error
}

// TokenRepository is an interface for interacting with token-related data.
type TokenRepository interface {
	// AddToken insets a new token into the database.
	AddToken(ctx context.Context, token *domain.Token) error
	// DeleteToken deletes a token with specified id.
	DeleteToken(ctx context.Context, id uuid.UUID) error
	// DeleteExpiredTokens deletes all tokens that have expired.
	DeleteExpiredTokens() error
}

// AuthService is an interface for interacting with auth-related business logic.
type AuthService interface {
	// Login validates user credentials and returns domain.TokenGroup.
	Login(ctx context.Context, user *domain.User) (*domain.TokenGroup, error)
	// RefreshSession uses a refresh token to refresh user session.
	RefreshSession(ctx context.Context, token *domain.Token) (*domain.TokenGroup, error)
}
