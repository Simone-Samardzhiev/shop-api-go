package mock

import (
	"context"
	"shop-api-go/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := ur.Called(ctx, user)
	return args.Error(0)
}
