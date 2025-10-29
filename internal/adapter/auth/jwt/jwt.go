package jwt

import (
	"shop-api-go/internal/adapter/config"
	"shop-api-go/internal/core/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// claims represent JWT token.
type claims struct {
	UserRole  domain.UserRole  `json:"userRole"`
	TokenType domain.TokenType `json:"tokenType"`
	jwt.RegisteredClaims
}

// TokenGenerator implements port.TokenGenerator and provides JWT generation.
type TokenGenerator struct {
	config *config.JWTConfig
}

// NewTokenGenerator creates a new TokenGenerator instance.
func NewTokenGenerator(config *config.JWTConfig) *TokenGenerator {
	return &TokenGenerator{config: config}
}

func (t *TokenGenerator) SignToken(token *domain.Token) (string, error) {
	now := time.Now()
	var exp time.Duration

	switch token.TokenType {
	case domain.AccessToken:
		exp = t.config.AccessTokenExpireTime
	case domain.RefreshToken:
		exp = t.config.RefreshTokenExpireTime
	default:
		zap.L().Error(
			"invalid token type",
			zap.String("token_type", string(token.TokenType)),
		)
		return "", domain.ErrInternal
	}

	token.ExpiresAt = now.Add(exp)
	jwtClaims := claims{
		UserRole:  token.UserRole,
		TokenType: token.TokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{t.config.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			ID:        token.Id.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    t.config.Issuer,
			NotBefore: jwt.NewNumericDate(now),
			Subject:   token.UserId.String(),
		},
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(t.config.Secret)
	if err != nil {
		zap.L().Error(
			"error signing token",
			zap.Error(err),
		)
		return "", domain.ErrInternal
	}
	return signedToken, nil
}

func (t *TokenGenerator) ParseToken(token string) (*domain.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return t.config.Secret, nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	jwtClaims, ok := parsedToken.Claims.(*claims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userId, err := uuid.Parse(jwtClaims.Subject)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	tokenId, err := uuid.Parse(jwtClaims.ID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	return domain.NewToken(tokenId, userId, jwtClaims.UserRole, jwtClaims.TokenType, jwtClaims.ExpiresAt.Time), nil
}
