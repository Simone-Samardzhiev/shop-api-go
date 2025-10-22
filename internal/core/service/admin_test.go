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

func TestAdminService_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		token          *domain.Token
		get            *domain.GetUsers
		expectedError  error
		expectedResult *domain.UsersResult
		mockSetup      func(
			mockUserRepository *mock.MockUserRepository,
			mockTokenRepository *mock.MockTokenRepository,
			mockPasswordHasher *mock.MockPasswordHasher,
		)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			get: &domain.GetUsers{
				Id: &uuid.UUID{},
			},
			expectedError: nil,
			expectedResult: &domain.UsersResult{
				Users: []domain.User{
					{
						Id:       uuid.UUID{},
						Username: "fetchedUser",
					},
				},
			},
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {
				mockUserRepository.
					EXPECT().
					GetUserById(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
					).
					Return(&domain.User{
						Id:       uuid.UUID{},
						Username: "fetchedUser",
					}, nil)

			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {

			},
		}, {
			name: "error invalid user role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
			},
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {

			},
		}, {
			name: "error limit not set",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			expectedError: domain.ErrLimitNotSet,
			get: &domain.GetUsers{
				After: &time.Time{},
			},
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {

			},
		}, {
			name: "error invalid query",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			get: &domain.GetUsers{
				Limit: new(int),
			},
			expectedError: domain.ErrInvalidQuery,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			tt.mockSetup(mockUserRepository, mockTokenRepository, mockPasswordHasher)

			result, err := service.
				NewAdminService(mockUserRepository, mockTokenRepository, mockPasswordHasher).
				GetUsers(context.Background(), tt.token, tt.get)

			if tt.expectedError == nil {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			} else {
				require.ErrorIs(t, err, tt.expectedError)
			}
		})
	}

}
func TestAdminService_UpdateUser(t *testing.T) {
	username := "newUsername"

	tests := []struct {
		name          string
		token         *domain.Token
		update        *domain.UserUpdate
		expectedError error
		mockSetup     func(
			mockUserRepository *mock.MockUserRepository,
			mockTokenRepository *mock.MockTokenRepository,
			mockPasswordHasher *mock.MockPasswordHasher,
		)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			update: &domain.UserUpdate{
				Id:       uuid.UUID{},
				Username: &username,
			},
			expectedError: nil,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {
				mockUserRepository.
					EXPECT().
					UpdateUser(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.Eq(&domain.UserUpdate{
							Id:       uuid.UUID{},
							Username: &username,
						}),
					).
					Return(nil)
				mockTokenRepository.
					EXPECT().
					DeleteAllTokensByUserId(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
					).
					Return(nil)
			},
		}, {
			name: "error invalid token type",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
				UserRole:  domain.Admin,
			},
			update: &domain.UserUpdate{
				Id:       uuid.UUID{},
				Username: &username,
			},
			expectedError: domain.ErrInvalidTokenType,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {

			},
		}, {
			name: "error invalid token role",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Client,
			},
			update: &domain.UserUpdate{
				Id:       uuid.UUID{},
				Username: &username,
			},
			expectedError: domain.ErrInvalidTokenRole,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {

			},
		}, {
			name: "error no fields to update",
			token: &domain.Token{
				TokenType: domain.AccessToken,
				UserRole:  domain.Admin,
			},
			update: &domain.UserUpdate{
				Id: uuid.UUID{},
			},
			expectedError: domain.ErrNoFieldsToUpdate,
			mockSetup: func(
				mockUserRepository *mock.MockUserRepository,
				mockTokenRepository *mock.MockTokenRepository,
				mockPasswordHasher *mock.MockPasswordHasher,
			) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)

			tt.mockSetup(mockUserRepository, mockTokenRepository, mockPasswordHasher)

			err := service.
				NewAdminService(mockUserRepository, mockTokenRepository, mockPasswordHasher).
				UpdateUser(context.Background(), tt.token, tt.update)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
