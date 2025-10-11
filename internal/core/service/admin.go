package service

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"time"

	"github.com/google/uuid"
)

// AdminService implements port.AdminService interface and provides access to admin-related business logic.
type AdminService struct {
	userRepository port.UserRepository
	passwordHasher port.PasswordHasher
}

// NewAdminService creates a new AdminService instance.
func NewAdminService(userRepository port.UserRepository, passwordHasher port.PasswordHasher) *AdminService {
	return &AdminService{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (s *AdminService) GetUsersByOffestPagination(ctx context.Context, token *domain.Token, page, limit int) ([]domain.User, error) {
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

func (s *AdminService) GetUsersByTimePagination(ctx context.Context, token *domain.Token, after time.Time, limit int) ([]domain.User, error) {
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

func (s *AdminService) SearchUserByUsername(ctx context.Context, token *domain.Token, username string, limit int) ([]domain.User, error) {
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

func (s *AdminService) SearchUserByEmail(ctx context.Context, token *domain.Token, email string, limit int) ([]domain.User, error) {
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

func (s *AdminService) GetUserById(ctx context.Context, token *domain.Token, id uuid.UUID) (*domain.User, error) {
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

func (s *AdminService) UpdateUsername(ctx context.Context, token *domain.Token, id uuid.UUID, username string) error {
	if token.TokenType != domain.AccessToken {
		return domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return domain.ErrInvalidTokenRole
	}

	return s.userRepository.UpdateUsername(ctx, id, username)
}

func (s *AdminService) UpdateEmail(ctx context.Context, token *domain.Token, id uuid.UUID, email string) error {
	if token.TokenType != domain.AccessToken {
		return domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return domain.ErrInvalidTokenRole
	}
	return s.userRepository.UpdateEmail(ctx, id, email)
}

func (s *AdminService) UpdatePassword(ctx context.Context, token *domain.Token, id uuid.UUID, password string) error {
	if token.TokenType != domain.AccessToken {
		return domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return domain.ErrInvalidTokenRole
	}

	hash, err := s.passwordHasher.Hash(password)
	if err != nil {
		return err
	}

	return s.userRepository.UpdatePassword(ctx, id, hash)
}

func (s *AdminService) UpdateRole(ctx context.Context, token *domain.Token, id uuid.UUID, role domain.UserRole) error {
	if token.TokenType != domain.AccessToken {
		return domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return domain.ErrInvalidTokenRole
	}

	return s.userRepository.UpdateRole(ctx, id, role)
}
