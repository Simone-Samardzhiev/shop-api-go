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

func TestUserService_ChangeUsername(t *testing.T) {
	tests := []struct {
		name          string
		user          *domain.User
		username      string
		expectedError error
		mockSetup     func(
			mockUserRepository *mock.MockUserRepository,
			mockPasswordHasher *mock.MockPasswordHasher,
			mockTokenRepository *mock.MockTokenRepository,
		)
	}{
		{
			name: "success",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			username:      "newUsername",
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
							"username",
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
						UpdateUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
							"newUsername",
						).
						Return(nil),
					mockTokenRepository.
						EXPECT().
						DeleteAllTokensByUserId(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
						).
						Return(nil),
				)
			},
		}, {
			name: "error wrong credentials",
			user: &domain.User{
				Username: "wrongUsername",
				Password: "password",
			},
			username:      "newUsername",
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
						"wrongUsername",
					).
					Return(nil, domain.ErrUserNotFound)
			},
		}, {
			name: "error fetching user",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			username:      "newUsername",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						"username",
					).
					Return(nil, domain.ErrInternalServerError)
			},
		}, {
			name: "error wrong credentials",
			user: &domain.User{
				Username: "username",
				Password: "wrongPassword",
			},
			username:      "newUsername",
			expectedError: domain.ErrWrongCredentials,
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
							"username",
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
			name: "error updating username",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			username:      "newUsername",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						"username",
					).
					Return(&domain.User{
						Username: "username",
						Password: "hashedPassword",
					}, nil)
				mockPasswordHasher.
					EXPECT().
					Compare("password", "hashedPassword").
					Return(nil)
				mockUserRepository.
					EXPECT().
					UpdateUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}), "newUsername").
					Return(domain.ErrInternalServerError)
			},
		}, {
			name: "error deleting tokens",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			username:      "newUsername",
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
							"username",
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
						UpdateUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
							"newUsername",
						).
						Return(nil),
					mockTokenRepository.
						EXPECT().
						DeleteAllTokensByUserId(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
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

			err := service.
				NewUserService(mockUserRepository, mockPasswordHasher, mockTokenRepository).
				ChangeUsername(context.Background(), tt.user, tt.username)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestUserService_ChangeEmail(t *testing.T) {
	tests := []struct {
		name          string
		user          *domain.User
		email         string
		expectedError error
		mockSetup     func(
			mockUserRepository *mock.MockUserRepository,
			mockPasswordHasher *mock.MockPasswordHasher,
			mockTokenRepository *mock.MockTokenRepository,
		)
	}{
		{
			name: "success",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			email:         "newEmail",
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
							"username",
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
						UpdateEmail(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
							"newEmail",
						).Return(nil),
				)
			},
		}, {
			name: "error fetching user",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			email:         "newEmail",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						"username",
					).
					Return(nil, domain.ErrInternalServerError)
			},
		}, {
			name: "error wrong credentials",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			email:         "newEmail",
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
						"username",
					).
					Return(nil, domain.ErrUserNotFound)
			},
		}, {
			name: "error wrong credentials",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			email:         "newEmail",
			expectedError: domain.ErrWrongCredentials,
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
							"username",
						).
						Return(&domain.User{
							Username: "username",
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("password", "hashedPassword").
						Return(domain.ErrWrongCredentials),
				)
			},
		}, {
			name: "error updating email",
			user: &domain.User{
				Username: "username",
				Password: "password",
			},
			email:         "newEmail",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						"username",
					).
					Return(&domain.User{
						Username: "username",
						Password: "hashedPassword",
					}, nil)
				mockPasswordHasher.
					EXPECT().
					Compare("password", "hashedPassword").
					Return(nil)
				mockUserRepository.
					EXPECT().
					UpdateEmail(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
						"newEmail",
					).
					Return(domain.ErrInternalServerError)
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

			err := service.
				NewUserService(mockUserRepository, mockPasswordHasher, mockTokenRepository).
				ChangeEmail(context.Background(), tt.user, tt.email)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	tests := []struct {
		name          string
		user          *domain.User
		newPassword   string
		expectedError error
		mockSetup     func(
			mockUserRepository *mock.MockUserRepository,
			mockPasswordHasher *mock.MockPasswordHasher,
			mockTokenRepository *mock.MockTokenRepository,
		)
	}{
		{
			name: "success",
			user: &domain.User{
				Username: "username",
				Password: "oldPassword",
			},
			newPassword:   "newPassword",
			expectedError: nil,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				fetchedUserID := uuid.New()
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(gomock.Any(), "username").
						Return(&domain.User{
							Id:       fetchedUserID,
							Username: "username",
							Password: "hashedOldPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("oldPassword", "hashedOldPassword").
						Return(nil),
					mockPasswordHasher.
						EXPECT().
						Hash("newPassword").
						Return("hashedNewPassword", nil),
					mockUserRepository.
						EXPECT().
						UpdatePassword(gomock.Any(), fetchedUserID, "hashedNewPassword").
						Return(nil),
					mockTokenRepository.
						EXPECT().
						DeleteAllTokensByUserId(gomock.Any(), fetchedUserID).
						Return(nil),
				)
			},
		},
		{
			name:          "error fetching user",
			user:          &domain.User{Username: "username", Password: "oldPassword"},
			newPassword:   "newPassword",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(gomock.Any(), "username").
					Return(nil, domain.ErrInternalServerError)
			},
		},
		{
			name:          "error wrong credentials (user not found)",
			user:          &domain.User{Username: "username", Password: "oldPassword"},
			newPassword:   "newPassword",
			expectedError: domain.ErrWrongCredentials,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(gomock.Any(), "username").
					Return(nil, domain.ErrUserNotFound)
			},
		},
		{
			name:          "error wrong password",
			user:          &domain.User{Username: "username", Password: "wrongPassword"},
			newPassword:   "newPassword",
			expectedError: domain.ErrWrongCredentials,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				fetchedUserID := uuid.New()
				mockUserRepository.
					EXPECT().
					GetUserByUsername(gomock.Any(), "username").
					Return(&domain.User{
						Id:       fetchedUserID,
						Username: "username",
						Password: "hashedOldPassword",
					}, nil)
				mockPasswordHasher.
					EXPECT().
					Compare("wrongPassword", "hashedOldPassword").
					Return(domain.ErrWrongCredentials)
			},
		},
		{
			name:          "error hashing new password",
			user:          &domain.User{Username: "username", Password: "oldPassword"},
			newPassword:   "newPassword",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				fetchedUserID := uuid.New()
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(gomock.Any(), "username").
						Return(&domain.User{
							Id:       fetchedUserID,
							Username: "username",
							Password: "hashedOldPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("oldPassword", "hashedOldPassword").
						Return(nil),
					mockPasswordHasher.
						EXPECT().
						Hash("newPassword").
						Return("", domain.ErrInternalServerError),
				)
			},
		},
		{
			name:          "error updating password",
			user:          &domain.User{Username: "username", Password: "oldPassword"},
			newPassword:   "newPassword",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				fetchedUserID := uuid.New()
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(gomock.Any(), "username").
						Return(&domain.User{
							Id:       fetchedUserID,
							Username: "username",
							Password: "hashedOldPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("oldPassword", "hashedOldPassword").
						Return(nil),
					mockPasswordHasher.
						EXPECT().
						Hash("newPassword").
						Return("hashedNewPassword", nil),
					mockUserRepository.
						EXPECT().
						UpdatePassword(gomock.Any(), fetchedUserID, "hashedNewPassword").
						Return(domain.ErrInternalServerError),
				)
			},
		},
		{
			name:          "error deleting tokens",
			user:          &domain.User{Username: "username", Password: "oldPassword"},
			newPassword:   "newPassword",
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
			) {
				fetchedUserID := uuid.New()
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(gomock.Any(), "username").
						Return(&domain.User{
							Id:       fetchedUserID,
							Username: "username",
							Password: "hashedOldPassword",
						}, nil),
					mockPasswordHasher.
						EXPECT().
						Compare("oldPassword", "hashedOldPassword").
						Return(nil),
					mockPasswordHasher.
						EXPECT().
						Hash("newPassword").
						Return("hashedNewPassword", nil),
					mockUserRepository.
						EXPECT().
						UpdatePassword(gomock.Any(), fetchedUserID, "hashedNewPassword").
						Return(nil),
					mockTokenRepository.
						EXPECT().
						DeleteAllTokensByUserId(gomock.Any(), fetchedUserID).
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

			err := service.
				NewUserService(mockUserRepository, mockPasswordHasher, mockTokenRepository).
				ChangePassword(context.Background(), tt.user, tt.newPassword)

			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
