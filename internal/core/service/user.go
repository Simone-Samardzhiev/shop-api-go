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

func (s *UserService) UpdateAccount(ctx context.Context, update *domain.UpdateAccount) error {
	hasFieldToUpdate := false
	if update.NewUsername != nil {
		hasFieldToUpdate = true
	}
	if update.NewPassword != nil {
		hasFieldToUpdate = true
	}
	if update.NewPassword != nil {
		hasFieldToUpdate = true
	}
	if !hasFieldToUpdate {
		return domain.ErrNoFieldsToUpdate
	}

	fetchedUser, err := s.userRepository.GetUserByUsername(ctx, update.Username)
	if errors.Is(err, domain.ErrUserNotFound) {
		return domain.ErrWrongCredentials
	} else if err != nil {
		return err
	}

	if err = s.passwordHasher.Compare(update.Password, fetchedUser.Password); err != nil {
		return err
	}

	if err = s.userRepository.UpdateUser(ctx, &domain.UserUpdate{
		Id:       fetchedUser.Id,
		Username: update.NewUsername,
		Email:    update.NewEmail,
		Password: update.NewPassword,
	}); err != nil {
		return err
	}

	if err = s.tokenRepository.DeleteAllTokensByUserId(ctx, fetchedUser.Id); err != nil {
		return err
	}

	return nil
}
