package service_test

import (
	"context"
	"errors"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name          string
		user          *domain.User
		expectedError error
		mockSetup     func(
			userRepository *mock.MockUserRepository,
			passwordHasher *mock.MockPasswordHasher,
			tokenRepository *mock.MockTokenRepository,
		)
	}{
		{
			name: "success",
			user: &domain.User{
				Password: "password",
			},
			expectedError: nil,
			mockSetup: func(
				userRepository *mock.MockUserRepository,
				passwordHasher *mock.MockPasswordHasher,
				tokenRepository *mock.MockTokenRepository,
			) {
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
			mockSetup: func(
				userRepository *mock.MockUserRepository,
				passwordHasher *mock.MockPasswordHasher,
				tokenRepository *mock.MockTokenRepository,
			) {
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
			mockSetup: func(
				userRepository *mock.MockUserRepository,
				passwordHasher *mock.MockPasswordHasher,
				tokenRepository *mock.MockTokenRepository,
			) {
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
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher, mockTokenRepository)

			err := service.NewUserService(mockUserRepository, mockPasswordHasher, mockTokenRepository).
				Register(context.Background(), tt.user)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestUserService_UpdateAccount(t *testing.T) {
	newUsername := "newUsername"

	tests := []struct {
		name          string
		update        *domain.UpdateAccount
		expectedError error
		mockSetup     func(
			mockUserRepository *mock.MockUserRepository,
			mockPasswordHasher *mock.MockPasswordHasher,
			mockTokenRepository *mock.MockTokenRepository,
		)
	}{
		{
			name: "success",
			update: &domain.UpdateAccount{
				Username:    "username",
				Password:    "password",
				NewUsername: &newUsername,
			},
			expectedError: nil,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq("username"),
						).
						Return(&domain.User{
							Username: "username",
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("password", "hashedPassword").
						Return(nil),
					mockUserRepository.
						EXPECT().
						UpdateUser(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq(&domain.UserUpdate{
								Username: &newUsername,
							}),
						).
						Return(nil),
					mockTokenRepository.
						EXPECT().
						DeleteAllTokensByUserId(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq(uuid.UUID{}),
						).
						Return(nil),
				)
			},
		}, {
			name:          "error no fields to update",
			update:        &domain.UpdateAccount{},
			expectedError: domain.ErrNoFieldsToUpdate,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {

			},
		}, {
			name: "error fetching user",
			update: &domain.UpdateAccount{
				Username:    "username",
				Password:    "password",
				NewUsername: &newUsername,
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq("username"),
						).
						Return(nil, domain.ErrInternalServerError),
				)
			},
		}, {
			name: "error wrong username",
			update: &domain.UpdateAccount{
				Username:    "wrongUsername",
				Password:    "password",
				NewUsername: &newUsername,
			},
			expectedError: domain.ErrWrongCredentials,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.Eq("wrongUsername"),
					).
					Return(nil, domain.ErrUserNotFound)
			},
		}, {
			name: "error wrong password",
			update: &domain.UpdateAccount{
				Username:    "username",
				Password:    "wrongPassword",
				NewUsername: &newUsername,
			},
			expectedError: domain.ErrWrongCredentials,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher, mockTokenRepository *mock.MockTokenRepository) {

				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq("username"),
						).
						Return(&domain.User{
							Username: "username",
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("wrongPassword", "hashedPassword").
						Return(domain.ErrWrongCredentials),
				)
			},
		}, {
			name: "error updating user",
			update: &domain.UpdateAccount{
				Username:    "username",
				Password:    "password",
				NewUsername: &newUsername,
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq("username"),
						).
						Return(&domain.User{
							Username: "username",
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("password", "hashedPassword").
						Return(nil),
					mockUserRepository.
						EXPECT().
						UpdateUser(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq(&domain.UserUpdate{
								Username: &newUsername,
							}),
						).
						Return(domain.ErrInternalServerError),
				)
			},
		}, {
			name: "error deleting tokens",
			update: &domain.UpdateAccount{
				Username:    "username",
				Password:    "password",
				NewUsername: &newUsername,
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq("username"),
						).
						Return(&domain.User{
							Username: "username",
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("password", "hashedPassword").
						Return(nil),
					mockUserRepository.
						EXPECT().
						UpdateUser(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq(&domain.UserUpdate{
								Username: &newUsername,
							}),
						).
						Return(nil),
					mockTokenRepository.
						EXPECT().
						DeleteAllTokensByUserId(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.Eq(uuid.UUID{}),
						).
						Return(domain.ErrInternalServerError),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher, mockTokenRepository)

			if err := service.
				NewUserService(mockUserRepository, mockPasswordHasher, mockTokenRepository).
				UpdateAccount(context.Background(), tt.update); err != nil {
				require.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}
