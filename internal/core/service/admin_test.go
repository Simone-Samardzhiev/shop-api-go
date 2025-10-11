package service_test

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAdminService_GetUsersByOffestPagination(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedUsers []domain.User
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: []domain.User{
				{
					Username: "fetchedUser",
				},
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					GetUsersByOffestPagination(
						gomock.AssignableToTypeOf(
							context.Background()),
						gomock.AssignableToTypeOf(0),
						gomock.AssignableToTypeOf(0)).
					Return([]domain.User{{
						Username: "fetchedUser",
					}}, nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Client,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error fetching users",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					GetUsersByOffestPagination(
						gomock.AssignableToTypeOf(
							context.Background()),
						gomock.AssignableToTypeOf(0),
						gomock.AssignableToTypeOf(0)).
					Return(nil, domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			passwordHasher := mock.NewMockPasswordHasher(ctrl)

			tt.mockSetup(mockUserRepository, passwordHasher)
			users, err := service.
				NewAdminService(mockUserRepository, passwordHasher).
				GetUsersByOffestPagination(
					context.Background(),
					tt.token,
					0,
					0,
				)

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedUsers, users)
		})
	}
}

func TestAdminService_GetUsersByOffestPagination_Error(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedUsers []domain.User
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: []domain.User{
				{
					Username: "fetchedUser",
				},
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					GetUsersByTimePagination(
						gomock.AssignableToTypeOf(
							context.Background()),
						gomock.AssignableToTypeOf(time.Time{}),
						gomock.AssignableToTypeOf(0)).
					Return([]domain.User{{
						Username: "fetchedUser",
					}}, nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Client,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error fetching users",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					GetUsersByTimePagination(
						gomock.AssignableToTypeOf(
							context.Background()),
						gomock.AssignableToTypeOf(time.Time{}),
						gomock.AssignableToTypeOf(0)).
					Return(nil, domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			passwordHasher := mock.NewMockPasswordHasher(ctrl)

			tt.mockSetup(mockUserRepository, passwordHasher)
			users, err := service.
				NewAdminService(mockUserRepository, passwordHasher).
				GetUsersByTimePagination(
					context.Background(),
					tt.token,
					time.Time{},
					0,
				)

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedUsers, users)
		})
	}

}

func TestAdminService_SearchUserByUsername(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedUsers []domain.User
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: []domain.User{
				{
					Username: "fetchedUser",
				},
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					SearchUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(""),
						gomock.AssignableToTypeOf(0)).
					Return([]domain.User{{
						Username: "fetchedUser",
					}}, nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Client,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error fetching users",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					SearchUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(""),
						gomock.AssignableToTypeOf(0)).
					Return(nil, domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher)

			users, err := service.
				NewAdminService(mockUserRepository, mockPasswordHasher).
				SearchUserByUsername(context.Background(), tt.token, "", 10)
			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedUsers, users)
		})
	}
}

func TestAdminService_SearchUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedUsers []domain.User
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: []domain.User{
				{
					Username: "fetchedUser",
				},
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					SearchUserByEmail(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(""),
						gomock.AssignableToTypeOf(0)).
					Return([]domain.User{{
						Username: "fetchedUser",
					}}, nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Client,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error fetching users",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUsers: nil,
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					SearchUserByEmail(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(""),
						gomock.AssignableToTypeOf(0)).
					Return(nil, domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher)

			users, err := service.
				NewAdminService(mockUserRepository, mockPasswordHasher).
				SearchUserByEmail(context.Background(), tt.token, "", 10)
			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedUsers, users)
		})
	}

}

func TestAdminService_GetUserById(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedUser  *domain.User
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUser: &domain.User{
				Username: "fetchedUser",
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					GetUserById(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
					).
					Return(&domain.User{
						Username: "fetchedUser",
					}, nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
				UserRole:  domain.Admin,
			},
			expectedUser:  nil,
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Client,
			},
			expectedUser:  nil,
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error fetching user",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedUser:  nil,
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					GetUserById(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
					).
					Return(nil, domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher)

			user, err := service.
				NewAdminService(
					mockUserRepository,
					mockPasswordHasher,
				).
				GetUserById(context.Background(), tt.token, uuid.UUID{})

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedUser, user)
		})
	}
}

func TestAdminService_UpdateUsername(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					UpdateUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
						gomock.AssignableToTypeOf(""),
					).Return(nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
			},
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error user not found",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrUserNotFound,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					UpdateUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
						gomock.AssignableToTypeOf(""),
					).Return(domain.ErrUserNotFound)
			},
		}, {
			name: "error updating user",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					UpdateUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
						gomock.AssignableToTypeOf(""),
					).Return(domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher)

			err := service.
				NewAdminService(
					mockUserRepository,
					mockPasswordHasher,
				).
				UpdateUsername(context.Background(), tt.token, uuid.UUID{}, "")
			require.ErrorIs(t, tt.expectedError, err)
		})
	}
}

func TestAdminService_UpdateEmail(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					UpdateEmail(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
						gomock.AssignableToTypeOf(""),
					).Return(nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
			},
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error user not found",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrUserNotFound,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					UpdateEmail(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
						gomock.AssignableToTypeOf(""),
					).Return(domain.ErrUserNotFound)
			},
		}, {
			name: "error updating user",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockUserRepository.
					EXPECT().
					UpdateEmail(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
						gomock.AssignableToTypeOf(""),
					).Return(domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher)

			err := service.
				NewAdminService(
					mockUserRepository,
					mockPasswordHasher,
				).
				UpdateEmail(context.Background(), tt.token, uuid.UUID{}, "")
			require.ErrorIs(t, tt.expectedError, err)
		})
	}
}

func TestAdminService_UpdatePassword(t *testing.T) {
	tests := []struct {
		name          string
		token         *domain.Token
		expectedError error
		mockSetup     func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: nil,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				gomock.InOrder(
					mockPasswordHasher.
						EXPECT().
						Hash(gomock.AssignableToTypeOf("")).
						Return("hashedPassword", nil),
					mockUserRepository.
						EXPECT().
						UpdatePassword(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
							gomock.AssignableToTypeOf(""),
						).Return(nil),
				)

			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
			},
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {

			},
		}, {
			name: "error hashing password",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				mockPasswordHasher.
					EXPECT().
					Hash(gomock.AssignableToTypeOf("")).
					Return("", domain.ErrInternalServerError)
			},
		}, {
			name: "error user not found",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrUserNotFound,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				gomock.InOrder(
					mockPasswordHasher.
						EXPECT().
						Hash(gomock.AssignableToTypeOf("")).
						Return("hashedPassword", nil),
					mockUserRepository.
						EXPECT().
						UpdatePassword(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
							gomock.AssignableToTypeOf(""),
						).Return(domain.ErrUserNotFound),
				)
			},
		}, {
			name: "error updating user",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInternalServerError,
			mockSetup: func(mockUserRepository *mock.MockUserRepository, mockPasswordHasher *mock.MockPasswordHasher) {
				gomock.InOrder(
					mockPasswordHasher.
						EXPECT().
						Hash(gomock.AssignableToTypeOf("")).
						Return("hashedPassword", nil),
					mockUserRepository.
						EXPECT().
						UpdatePassword(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{}),
							gomock.AssignableToTypeOf(""),
						).Return(domain.ErrInternalServerError),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockPasswordHasher)

			err := service.
				NewAdminService(
					mockUserRepository,
					mockPasswordHasher,
				).
				UpdatePassword(context.Background(), tt.token, uuid.UUID{}, "")
			require.ErrorIs(t, tt.expectedError, err)
		})
	}
}
