package service

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

func TestAuthService_Login(t *testing.T) {
	mockTokenGenerator := new(mock.TokenGenerator)
	mockTokenRepository := new(mock.TokenRepository)
	mockUserRepository := new(mock.UserRepository)

	id := uuid.New()
	hash, err := util.HashPassword("password")
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	mockTokenGenerator.On("SignToken", tmock.Anything).Return("token", nil)

	mockTokenRepository.On("AddToken", tmock.Anything, tmock.Anything).Return(nil)

	mockUserRepository.On("GetUserByUsername", tmock.Anything, tmock.MatchedBy(func(username string) bool {
		return username == "NewUser"
	})).Return((*domain.User)(nil), domain.ErrWrongCredentials)
	mockUserRepository.On("GetUserByUsername", tmock.Anything, tmock.MatchedBy(func(username string) bool {
		return username == "ExistingUser"
	})).Return(&domain.User{
		Id:        id,
		Username:  "ExistingUser",
		Email:     "ExistingUser",
		Password:  string(hash),
		Role:      domain.Client,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	service := NewAuthService(mockTokenGenerator, mockTokenRepository, mockUserRepository)
	tests := []struct {
		name               string
		user               *domain.User
		expectedTokenGroup *domain.TokenGroup
		expectedErr        error
	}{
		{
			name:               "Success",
			user:               &domain.User{Username: "ExistingUser", Password: "password"},
			expectedTokenGroup: &domain.TokenGroup{AccessToken: "token", RefreshToken: "token"},
			expectedErr:        nil,
		},
		{
			name:               "Error invalid username",
			user:               &domain.User{Username: "NewUser", Password: "password"},
			expectedTokenGroup: nil,
			expectedErr:        domain.ErrWrongCredentials,
		}, {
			name:               "Error invalid password",
			user:               &domain.User{Username: "ExistingUser", Password: "pass"},
			expectedTokenGroup: nil,
			expectedErr:        domain.ErrWrongCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, serviceErr := service.Login(context.Background(), tt.user)
			if tt.expectedErr == nil {
				assert.Equal(t, res, tt.expectedTokenGroup)
			} else {
				assert.ErrorIs(t, serviceErr, tt.expectedErr)
			}
		})
	}
}
