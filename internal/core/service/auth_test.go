package service_test

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/service"
	"shop-api-go/internal/core/util"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTokenGenerator := mock.NewMockTokenGenerator(ctrl)
	mockTokenRepository := mock.NewMockTokenRepository(ctrl)
	mockUserRepository := mock.NewMockUserRepository(ctrl)

	hash, err := util.HashPassword("password")
	assert.NoError(t, err)

	gomock.InOrder(
		// Test 1 (success)
		mockUserRepository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).
			Return(&domain.User{Username: "existingUser", Password: string(hash)}, nil),
		mockTokenGenerator.EXPECT().
			SignToken(gomock.Any()).
			Return("token", nil).
			Times(2),
		mockTokenRepository.EXPECT().
			AddToken(gomock.Any(), gomock.Any()).
			Return(nil),

		// Test 2 (error fetching user)
		mockUserRepository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).
			Return((*domain.User)(nil), domain.ErrInternalServerError),

		// Test 3 (username is not found)
		mockUserRepository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).
			Return((*domain.User)(nil), domain.ErrWrongCredentials),

		// Test 4 (error signing token)
		mockUserRepository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).
			Return(&domain.User{Username: "existingUser", Password: string(hash)}, nil),
		mockTokenGenerator.EXPECT().
			SignToken(gomock.Any()).
			Return("", domain.ErrInternalServerError),

		// Test 5 (error adding token)
		mockUserRepository.EXPECT().
			GetUserByUsername(gomock.Any(), gomock.Any()).
			Return(&domain.User{Username: "existingUser", Password: string(hash)}, nil),
		mockTokenGenerator.EXPECT().
			SignToken(gomock.Any()).
			Return("token", nil).
			Times(2),
		mockTokenRepository.EXPECT().
			AddToken(gomock.Any(), gomock.Any()).
			Return(domain.ErrInternalServerError),
	)

	s := service.NewAuthService(mockTokenGenerator, mockTokenRepository, mockUserRepository)

	tests := []struct {
		name               string
		user               *domain.User
		expectedErr        error
		expectedTokenGroup *domain.TokenGroup
	}{
		{
			name:               "success",
			user:               &domain.User{Username: "existingUsername", Password: "password"},
			expectedErr:        nil,
			expectedTokenGroup: &domain.TokenGroup{AccessToken: "token", RefreshToken: "token"},
		}, {
			name:               "error fetching user",
			user:               &domain.User{},
			expectedErr:        domain.ErrInternalServerError,
			expectedTokenGroup: nil,
		},
		{
			name:               "username is not found",
			user:               &domain.User{},
			expectedErr:        domain.ErrWrongCredentials,
			expectedTokenGroup: nil,
		}, {
			name:               "error signing token",
			user:               &domain.User{Username: "existingUser", Password: "password"},
			expectedErr:        domain.ErrInternalServerError,
			expectedTokenGroup: nil,
		}, {
			name:               "error adding token",
			user:               &domain.User{Username: "existingUser", Password: "password"},
			expectedErr:        domain.ErrInternalServerError,
			expectedTokenGroup: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenGroup, serviceErr := s.Login(context.Background(), tt.user)
			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedTokenGroup, tokenGroup)
			} else {
				assert.ErrorIs(t, tt.expectedErr, serviceErr)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockTokenGenerator := mock.NewMockTokenGenerator(gomock.NewController(t))
	mockTokenRepository := mock.NewMockTokenRepository(gomock.NewController(t))
	mockUserRepository := mock.NewMockUserRepository(gomock.NewController(t))

	gomock.InOrder(
		// Test 1(success)
		mockTokenRepository.EXPECT().
			DeleteToken(gomock.Any(), gomock.AssignableToTypeOf(uuid.UUID{})).
			Return(true, nil),
		mockTokenGenerator.EXPECT().
			SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
			Return("token", nil).
			Times(2),
		mockTokenRepository.EXPECT().
			AddToken(gomock.Any(), gomock.AssignableToTypeOf(&domain.Token{})).
			Return(nil),

		// Test 2(invalid token)
		mockTokenRepository.EXPECT().
			DeleteToken(gomock.Any(), gomock.AssignableToTypeOf(uuid.UUID{})).
			Return(false, nil),

		// Test 3(error deleting token)
		mockTokenRepository.EXPECT().
			DeleteToken(gomock.Any(), gomock.AssignableToTypeOf(uuid.UUID{})).
			Return(false, domain.ErrInternalServerError),

		// Test 4(error signing token)
		mockTokenRepository.EXPECT().
			DeleteToken(gomock.Any(), gomock.AssignableToTypeOf(uuid.UUID{})).
			Return(true, nil),
		mockTokenGenerator.EXPECT().
			SignToken(gomock.AssignableToTypeOf(&domain.Token{})).Return("", domain.ErrInternalServerError),

		// Test 5(error adding token)
		mockTokenRepository.EXPECT().
			DeleteToken(gomock.Any(), gomock.AssignableToTypeOf(uuid.UUID{})).
			Return(true, nil),
		mockTokenGenerator.EXPECT().
			SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
			Return("token", nil).
			Times(2),
		mockTokenRepository.EXPECT().AddToken(gomock.Any(), gomock.AssignableToTypeOf(&domain.Token{})).
			Return(domain.ErrInternalServerError),
	)

	s := service.NewAuthService(mockTokenGenerator, mockTokenRepository, mockUserRepository)

	tests := []struct {
		name               string
		token              *domain.Token
		expectedTokenGroup *domain.TokenGroup
		expectedErr        error
	}{
		{
			name:               "success",
			token:              &domain.Token{},
			expectedTokenGroup: &domain.TokenGroup{AccessToken: "token", RefreshToken: "token"},
			expectedErr:        nil,
		}, {
			name:               "invalid token",
			token:              &domain.Token{},
			expectedTokenGroup: nil,
			expectedErr:        domain.ErrInvalidToken,
		}, {
			name:               "error deleting token",
			token:              &domain.Token{},
			expectedTokenGroup: nil,
			expectedErr:        domain.ErrInternalServerError,
		}, {
			name:               "error signing token",
			token:              &domain.Token{},
			expectedTokenGroup: nil,
			expectedErr:        domain.ErrInternalServerError,
		}, {
			name:               "error adding token",
			token:              &domain.Token{},
			expectedTokenGroup: nil,
			expectedErr:        domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenGroup, serviceErr := s.RefreshSession(context.Background(), tt.token)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedTokenGroup, tokenGroup)
			} else {
				assert.ErrorIs(t, tt.expectedErr, serviceErr)
			}
		})
	}
}
