package service

import (
	"context"
	"errors"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"testing"

	temock "github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

func TestUserService_Register(t *testing.T) {
	mockRepo := new(mock.UserRepository)
	ctx := context.Background()
	id := uuid.New()
	// Success case
	mockRepo.On("CreateUser", ctx, temock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "email" && u.Username == "username"
	})).Return(nil)

	// Duplicate email case
	mockRepo.On("CreateUser", ctx, temock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "duplicate"
	})).Return(domain.ErrEmailAlreadyInUse)

	// Duplicate username case
	mockRepo.On("CreateUser", ctx, temock.MatchedBy(func(u *domain.User) bool {
		return u.Username == "duplicate"
	})).Return(domain.ErrUsernameAlreadyInUse)

	service := NewUserService(mockRepo)

	tests := []struct {
		name          string
		user          *domain.User
		expectedError error
	}{
		{
			name: "Success",
			user: &domain.User{
				Id:       id,
				Email:    "email",
				Username: "username",
				Password: "password",
				Role:     domain.Client,
			},
			expectedError: nil,
		}, {
			name: "Error with duplicate email",
			user: &domain.User{
				Id:       id,
				Email:    "duplicate",
				Username: "username",
				Password: "password",
				Role:     domain.Client,
			},
			expectedError: domain.ErrEmailAlreadyInUse,
		}, {
			name: "Error with duplicate username",
			user: &domain.User{
				Id:       id,
				Email:    "email",
				Username: "duplicate",
				Password: "password",
				Role:     domain.Client,
			},
			expectedError: domain.ErrUsernameAlreadyInUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Register(ctx, tt.user)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}
