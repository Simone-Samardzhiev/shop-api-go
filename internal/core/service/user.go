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
	userRepository port.UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(ur port.UserRepository) *UserService {
	return &UserService{userRepository: ur}
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
	return s.userRepository.AddUser(ctx, user)
}

func (s *UserService) GetUsersByOffestPagination(ctx context.Context, token *domain.Token, page, limit int) ([]domain.User, error) {
	if token.TokenType != domain.AccessToken {
		return nil, domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return nil, domain.ErrInvalidTokenRole
	}

	users, err := s.userRepository.GetUsersByOffestPagination(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUsersByTimePagination(ctx context.Context, token *domain.Token, after time.Time, limit int) ([]domain.User, error) {
	if token.TokenType != domain.AccessToken {
		return nil, domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return nil, domain.ErrInvalidTokenRole
	}

	users, err := s.userRepository.GetUsersByTimePagination(ctx, after, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) SearchUserByUsername(ctx context.Context, token *domain.Token, username string, limit int) ([]domain.User, error) {
	if token.TokenType != domain.AccessToken {
		return nil, domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return nil, domain.ErrInvalidTokenRole
	}

	users, err := s.userRepository.SearchUserByUsername(ctx, username, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) SearchUserByEmail(ctx context.Context, token *domain.Token, email string, limit int) ([]domain.User, error) {
	if token.TokenType != domain.AccessToken {
		return nil, domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return nil, domain.ErrInvalidTokenRole
	}

	users, err := s.userRepository.SearchUserByEmail(ctx, email, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserById(ctx context.Context, token *domain.Token, id uuid.UUID) (*domain.User, error) {
	if token.TokenType != domain.AccessToken {
		return nil, domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return nil, domain.ErrInvalidTokenRole
	}

	user, err := s.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUsername(ctx context.Context, token *domain.Token, id uuid.UUID, username string) error {
	if token.TokenType != domain.AccessToken {
		return domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return domain.ErrInvalidTokenRole
	}

	return s.userRepository.UpdateUsername(ctx, id, username)
}
