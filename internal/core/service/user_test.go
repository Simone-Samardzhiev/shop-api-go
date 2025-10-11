package service_test

import (
	"context"
	"errors"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/service"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name          string
		user          *domain.User
		expectedError error
		mockSetup     func(userRepository *mock.MockUserRepository, passwordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			user: &domain.User{
				Password: "password",
			},
			expectedError: nil,
			mockSetup: func(userRepository *mock.MockUserRepository, passwordHasher *mock.MockPasswordHasher) {
				gomock.InOrder(
					passwordHasher.EXPECT().
						Hash("password").
						Return("hashedPassword", nil),

					userRepository.EXPECT().
						AddUser(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(&domain.User{}),
						).
						DoAndReturn(func(ctx context.Context, user *domain.User) error {
							if user.Password != "hashedPassword" {
								return errors.New("wrong password")
							}
							return nil
						}),
				)
			},
		}, {
			name: "error hashing password",
			user: &domain.User{
				Password: "password",
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(userRepository *mock.MockUserRepository, passwordHasher *mock.MockPasswordHasher) {
				passwordHasher.EXPECT().
					Hash("password").
					Return("", domain.ErrInternalServerError)
			},
		}, {
			name: "error adding user",
			user: &domain.User{
				Password: "password",
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(userRepository *mock.MockUserRepository, passwordHasher *mock.MockPasswordHasher) {
				gomock.InOrder(
					passwordHasher.EXPECT().
						Hash("password").
						Return("hashedPassword", nil),
					userRepository.EXPECT().
						AddUser(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(&domain.User{}),
						).
						DoAndReturn(func(ctx context.Context, user *domain.User) error {
							if user.Password != "hashedPassword" {
								return errors.New("wrong password")
							}
							return domain.ErrInternalServerError
						}),
				)

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher)

			err := service.NewUserService(mockUserRepository, mockPasswordHasher).Register(context.Background(), tt.user)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
