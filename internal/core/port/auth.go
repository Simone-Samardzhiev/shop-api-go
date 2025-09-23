package port

import (
	"context"
	"shop-api-go/internal/core/domain"
)

// TokenGenerator is an interface for signing and decoding tokens.
type TokenGenerator interface {
	// SignToken returns a signed string token.
	SignToken(token *domain.Token) (string, error)
	// ParseToken takes a signed token and returns *domain.Token
	ParseToken(token string) (*domain.Token, error)
}

// TokenRepository is an interface for interacting with token-related data.
type TokenRepository interface {
	// AddToken insets a new token into the database.
	AddToken(ctx context.Context, token *domain.Token) error
}

// AuthService is an interface for interacting with auth-related business logic.
type AuthService interface {
	// Login validates user credentials and returns domain.TokenGroup.
	Login(ctx context.Context, user *domain.User) (*domain.TokenGroup, error)
}
