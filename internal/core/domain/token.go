package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenType is an enum for token's type.
type TokenType string

// TokenType enum values.
const (
	AccessToken  = TokenType("access_token")
	RefreshToken = TokenType("refresh_token")
)

// Token is an entity representing a token.
type Token struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	TokenType TokenType
	UserRole  UserRole
	ExpiresAt time.Time
}

// TokenGroup is an entity representing a group of signed access and refresh token.
type TokenGroup struct {
	AccessToken  string
	RefreshToken string
}
