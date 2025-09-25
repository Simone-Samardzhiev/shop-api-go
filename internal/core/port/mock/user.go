package mock

import (
	"context"
	"shop-api-go/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (ur *UserRepository) AddUser(ctx context.Context, user *domain.User) error {
	args := ur.Called(ctx, user)
	return args.Error(0)
}

func (ur *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := ur.Called(ctx, username)
	return args.Get(0).(*domain.User), args.Error(1)
}
