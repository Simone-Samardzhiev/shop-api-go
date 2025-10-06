package service_test

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAdminService_GetUsersByPages(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	gomock.InOrder(
		mockUserRepository.EXPECT().
			GetUsersByOffestPagination(
				gomock.Any(),
				gomock.AssignableToTypeOf(0),
				gomock.AssignableToTypeOf(0),
			).
			Return([]domain.User{
				{
					Username: "fetchedUser",
				},
			}, nil),
		mockUserRepository.EXPECT().
			GetUsersByOffestPagination(
				gomock.Any(),
				gomock.AssignableToTypeOf(0),
				gomock.AssignableToTypeOf(0),
			).
			Return(nil, domain.ErrInternalServerError),
	)

	s := service.NewAdminService(mockUserRepository)
	tests := []struct {
		name          string
		expectedUsers []domain.User
		expectedErr   error
	}{
		{
			name: "success",
			expectedUsers: []domain.User{
				{
					Username: "fetchedUser",
				},
			},
			expectedErr: nil,
		}, {
			name:          "error",
			expectedUsers: nil,
			expectedErr:   domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetchedUsers, serviceErr := s.GetUsersByOffestPagination(
				context.Background(),
				&domain.Token{
					TokenType: domain.AccessToken,
					UserRole:  domain.Admin,
				},
				0,
				0,
			)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedUsers, fetchedUsers)
			} else {
				assert.ErrorIs(t, tt.expectedErr, serviceErr)
			}
		})
	}
}

func TestAdminService_GetUsersByTimePagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	gomock.InOrder(
		mockUserRepository.EXPECT().
			GetUsersByTimePagination(
				gomock.Any(),
				gomock.AssignableToTypeOf(time.Time{}),
				gomock.AssignableToTypeOf(0),
			).
			Return([]domain.User{
				{
					Username: "fetchedUser",
				},
			}, nil),
		mockUserRepository.EXPECT().
			GetUsersByTimePagination(
				gomock.Any(),
				gomock.AssignableToTypeOf(time.Time{}),
				gomock.AssignableToTypeOf(0),
			).
			Return(nil, domain.ErrInternalServerError),
	)

	s := service.NewAdminService(mockUserRepository)
	tests := []struct {
		name          string
		expectedUsers []domain.User
		expectedErr   error
	}{
		{
			name: "success",
			expectedUsers: []domain.User{{
				Username: "fetchedUser",
			}},
			expectedErr: nil,
		},
		{
			name:          "error",
			expectedUsers: nil,
			expectedErr:   domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, serviceErr := s.GetUsersByTimePagination(
				context.Background(),
				&domain.Token{
					TokenType: domain.AccessToken,
					UserRole:  domain.Admin,
				},
				time.Time{},
				0,
			)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedUsers, users)
			} else {
				assert.ErrorIs(t, tt.expectedErr, serviceErr)
			}
		})
	}
}

func TestAdminService_SearchUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	gomock.InOrder(
		mockUserRepository.EXPECT().
			SearchUserByUsername(
				gomock.Any(),
				gomock.AssignableToTypeOf(""),
				gomock.AssignableToTypeOf(0),
			).
			Return([]domain.User{
				{
					Username: "fetchedUser"},
			}, nil),
		mockUserRepository.EXPECT().
			SearchUserByUsername(gomock.Any(),
				gomock.AssignableToTypeOf(""),
				gomock.AssignableToTypeOf(0),
			).
			Return(nil, domain.ErrInternalServerError),
	)

	s := service.NewAdminService(mockUserRepository)
	tests := []struct {
		name          string
		expectedUsers []domain.User
		expectedErr   error
	}{
		{
			name: "success",
			expectedUsers: []domain.User{
				{
					Username: "fetchedUser",
				},
			},
			expectedErr: nil,
		}, {
			name:          "error",
			expectedUsers: nil,
			expectedErr:   domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, serviceErr := s.SearchUserByUsername(
				context.Background(),
				&domain.Token{
					TokenType: domain.AccessToken,
					UserRole:  domain.Admin,
				},
				"",
				1)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedUsers, users)
			} else {
				assert.ErrorIs(t, tt.expectedErr, serviceErr)
			}
		})
	}
}

func TestAdminService_SearchUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	gomock.InOrder(
		mockUserRepository.EXPECT().
			SearchUserByEmail(
				gomock.Any(),
				gomock.AssignableToTypeOf(""),
				gomock.AssignableToTypeOf(0),
			).
			Return([]domain.User{
				{
					Username: "fetchedUser"},
			}, nil),
		mockUserRepository.EXPECT().
			SearchUserByEmail(gomock.Any(),
				gomock.AssignableToTypeOf(""),
				gomock.AssignableToTypeOf(0),
			).
			Return(nil, domain.ErrInternalServerError),
	)

	s := service.NewAdminService(mockUserRepository)
	tests := []struct {
		name          string
		expectedUsers []domain.User
		expectedErr   error
	}{
		{
			name: "success",
			expectedUsers: []domain.User{
				{
					Username: "fetchedUser",
				},
			},
			expectedErr: nil,
		}, {
			name:          "error",
			expectedUsers: nil,
			expectedErr:   domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, serviceErr := s.SearchUserByEmail(
				context.Background(),
				&domain.Token{
					TokenType: domain.AccessToken,
					UserRole:  domain.Admin,
				},
				"",
				1)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedUsers, users)
			} else {
				assert.ErrorIs(t, tt.expectedErr, serviceErr)
			}
		})
	}
}

func TestAdminService_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	gomock.InOrder(
		mockUserRepository.EXPECT().
			GetUserById(
				gomock.Any(),
				gomock.AssignableToTypeOf(uuid.UUID{})).
			Return(&domain.User{
				Username: "fetchedUser",
			}, nil),
		mockUserRepository.EXPECT().
			GetUserById(
				gomock.Any(),
				gomock.AssignableToTypeOf(uuid.UUID{})).
			Return(nil, domain.ErrInternalServerError),
	)

	s := service.NewAdminService(mockUserRepository)

	tests := []struct {
		name         string
		expectedUser *domain.User
		expectedErr  error
	}{
		{
			name: "success",
			expectedUser: &domain.User{
				Username: "fetchedUser",
			},
			expectedErr: nil,
		}, {
			name:         "error",
			expectedUser: nil,
			expectedErr:  domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := s.GetUserById(
				context.Background(),
				&domain.Token{
					TokenType: domain.AccessToken,
					UserRole:  domain.Admin,
				},
				uuid.UUID{},
			)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedUser, user)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestAdminService_UpdateUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	gomock.InOrder(
		mockUserRepository.EXPECT().
			UpdateUsername(
				gomock.Any(),
				gomock.AssignableToTypeOf(uuid.UUID{}),
				gomock.AssignableToTypeOf("")).
			Return(nil),
		mockUserRepository.EXPECT().
			UpdateUsername(
				gomock.Any(),
				gomock.AssignableToTypeOf(uuid.UUID{}),
				gomock.AssignableToTypeOf("")).
			Return(domain.ErrUserNotFound),
		mockUserRepository.EXPECT().
			UpdateUsername(
				gomock.Any(),
				gomock.AssignableToTypeOf(uuid.UUID{}),
				gomock.AssignableToTypeOf("")).
			Return(domain.ErrInternalServerError),
	)

	s := service.NewAdminService(mockUserRepository)

	tests := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "success",
			expectedErr: nil,
		}, {
			name:        "not found",
			expectedErr: domain.ErrUserNotFound,
		}, {
			name:        "error",
			expectedErr: domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceErr := s.UpdateUsername(
				context.Background(),
				&domain.Token{
					TokenType: domain.AccessToken,
					UserRole:  domain.Admin,
				},
				uuid.UUID{},
				"",
			)

			assert.ErrorIs(t, serviceErr, tt.expectedErr)
		})
	}
}
