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

// NewToken creates a new Token instance.
func NewToken(id, userId uuid.UUID, userRole UserRole, tokenType TokenType, expires time.Time) *Token {
	return &Token{
		Id:        id,
		UserId:    userId,
		TokenType: tokenType,
		UserRole:  userRole,
		ExpiresAt: expires,
	}
}

// TokenGroup is an entity representing a group of signed access and refresh token.
type TokenGroup struct {
	AccessToken  string
	RefreshToken string
}
