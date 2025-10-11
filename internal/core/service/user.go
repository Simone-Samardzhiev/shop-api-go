package service

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"time"

	"github.com/google/uuid"
)

// UserService implements port.UserService interface and provides access to user-related business logic.
type UserService struct {
	userRepository port.UserRepository
	passwordHasher port.PasswordHasher
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepository port.UserRepository, passwordHasher port.PasswordHasher) *UserService {
	return &UserService{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	user.Id = uuid.New()
	hash, err := s.passwordHasher.Hash(user.Password)
	if err != nil {
		return domain.ErrInternalServerError
	}
	user.Password = hash
	return s.userRepository.AddUser(ctx, user)
}
