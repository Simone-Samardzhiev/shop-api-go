package service

import (
	"context"
	"errors"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"

	"github.com/google/uuid"
)

// UserService implements port.UserService interface and provides access to user-related business logic.
type UserService struct {
	userRepository  port.UserRepository
	passwordHasher  port.PasswordHasher
	tokenRepository port.TokenRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(
	userRepository port.UserRepository,
	passwordHasher port.PasswordHasher,
	tokenRepository port.TokenRepository,
) *UserService {
	return &UserService{
		userRepository:  userRepository,
		passwordHasher:  passwordHasher,
		tokenRepository: tokenRepository,
	}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) error {
	user.Id = uuid.New()
	hash, err := s.passwordHasher.Hash(user.Password)
	if err != nil {
		return domain.ErrInternalServerError
	}
	user.Password = hash
	return s.userRepository.AddUser(ctx, user)
}

func (s *UserService) ChangeUsername(ctx context.Context, user *domain.User, username string) error {
	fetchedUser, err := s.userRepository.GetUserByUsername(ctx, user.Username)
	if errors.Is(err, domain.ErrUserNotFound) {
		return domain.ErrWrongCredentials
	} else if err != nil {
		return domain.ErrInternalServerError
	}

	err = s.passwordHasher.Compare(user.Password, fetchedUser.Password)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdateUsername(ctx, fetchedUser.Id, username)
	if err != nil {
		return domain.ErrInternalServerError
	}

	err = s.tokenRepository.DeleteAllTokensByUserId(ctx, fetchedUser.Id)
	if err != nil {
		return domain.ErrInternalServerError
	}

	return nil
}
