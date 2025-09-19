package service

import (
	"context"
	"shop-api-go/internal/adapter/storage/repository"
	"shop-api-go/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService implements port.UserService interface and provides access to user-related business logic.
type UserService struct {
	ur repository.UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(ur repository.UserRepository) *UserService {
	return &UserService{ur: ur}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	user.Id = uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)

	err = s.ur.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
