package service

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"shop-api-go/internal/core/util"
	"time"

	"github.com/google/uuid"
)

// UserService implements port.UserService interface and provides access to user-related business logic.
type UserService struct {
	ur port.UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(ur port.UserRepository) *UserService {
	return &UserService{ur: ur}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	user.Id = uuid.New()
	hash, err := util.HashPassword(user.Password)
	if err != nil {
		return domain.ErrInternalServerError
	}
	user.Password = string(hash)
	return s.ur.AddUser(ctx, user)
}

func (s *UserService) GetUsersByOffestPagination(ctx context.Context, token *domain.Token, page, limit int) ([]domain.User, error) {
	if token.TokenType != domain.AccessToken {
		return nil, domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return nil, domain.ErrInvalidTokenRole
	}

	users, err := s.ur.GetUsersByOffestPagination(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}
