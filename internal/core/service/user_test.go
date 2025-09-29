package service_test

import (
	"context"
	"errors"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUserRepository := mock.NewMockUserRepository(ctrl)
	gomock.InOrder(
		mockUserRepository.EXPECT().
			AddUser(gomock.Any(), gomock.AssignableToTypeOf(&domain.User{})).
			DoAndReturn(func(ctx context.Context, user *domain.User) error {
				return nil
			}),
		mockUserRepository.EXPECT().
			AddUser(gomock.Any(), gomock.AssignableToTypeOf(&domain.User{})).
			DoAndReturn(func(ctx context.Context, user *domain.User) error {
				if user.Username != "duplicate" {
					return errors.New("username is not duplicate")
				}
				return domain.ErrUsernameAlreadyInUse
			}),
		mockUserRepository.EXPECT().
			AddUser(gomock.Any(), gomock.AssignableToTypeOf(&domain.User{})).
			DoAndReturn(func(ctx context.Context, user *domain.User) error {
				if user.Email != "duplicate" {
					return errors.New("email is not duplicate")
				}
				return domain.ErrEmailAlreadyInUse
			}),
	)

	s := service.NewUserService(mockUserRepository)

	tests := []struct {
		name        string
		user        *domain.User
		expectedErr error
	}{
		{
			name:        "success",
			user:        &domain.User{},
			expectedErr: nil,
		},
		{
			name: "duplicate username",
			user: &domain.User{
				Username: "duplicate",
			},
			expectedErr: domain.ErrUsernameAlreadyInUse,
		}, {
			name: "duplicate email",
			user: &domain.User{
				Email: "duplicate",
			},
			expectedErr: domain.ErrEmailAlreadyInUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Register(context.Background(), tt.user)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
