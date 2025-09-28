package mock

import (
	"context"
	"shop-api-go/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type TokenGenerator struct {
	mock.Mock
}

func (m *TokenGenerator) SignToken(token *domain.Token) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *TokenGenerator) ParseToken(token string) (*domain.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.Token), args.Error(1)
}

type TokenRepository struct {
	mock.Mock
}

func (t *TokenRepository) AddToken(ctx context.Context, token *domain.Token) error {
	args := t.Called(ctx, token)
	return args.Error(0)
}

func (t *TokenRepository) DeleteExpiredTokens() error {
	args := t.Called()
	return args.Error(0)
}
