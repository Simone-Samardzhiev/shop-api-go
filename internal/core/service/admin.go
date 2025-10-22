package service

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"time"
)

// AdminService implements port.AdminService interface and provides access to admin-related business logic.
type AdminService struct {
	userRepository  port.UserRepository
	tokenRepository port.TokenRepository
	passwordHasher  port.PasswordHasher
}

// NewAdminService creates a new AdminService instance.
func NewAdminService(userRepository port.UserRepository, tokenRepository port.TokenRepository, passwordHasher port.PasswordHasher) *AdminService {
	return &AdminService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		passwordHasher:  passwordHasher,
	}
}

func (s *AdminService) GetUsers(ctx context.Context, token *domain.Token, get *domain.GetUsers) (*domain.UsersResult, error) {
	if token.TokenType != domain.AccessToken {
		return nil, domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return nil, domain.ErrInvalidTokenRole
	}
	if get.Id == nil && get.Limit == nil {
		return nil, domain.ErrLimitNotSet
	}

	switch {
	case get.Id != nil:
		user, err := s.userRepository.GetUserById(ctx, *get.Id)
		if err != nil {
			return nil, err
		}

		return domain.NewUsersResult([]domain.User{*user}, nil), nil
	case get.Username != nil:
		users, err := s.userRepository.SearchUserByUsername(ctx, *get.Username, *get.Limit, get.Role)
		if err != nil {
			return nil, err
		}

		return domain.NewUsersResult(users, nil), nil
	case get.Email != nil:
		users, err := s.userRepository.SearchUserByEmail(ctx, *get.Email, *get.Limit, get.Role)
		if err != nil {
			return nil, err
		}

		return domain.NewUsersResult(users, nil), nil
	case get.Page != nil:
		users, err := s.userRepository.GetUsersByOffestPagination(ctx, *get.Page, *get.Limit, get.Role)
		if err != nil {
			return nil, err
		}

		return domain.NewUsersResult(users, nil), nil
	case get.After != nil:
		users, err := s.userRepository.GetUsersByTimePagination(ctx, *get.After, *get.Limit, get.Role)
		if err != nil {
			return nil, err
		}

		var cursor string
		if len(users) > 0 {
			cursor = users[len(users)-1].CreatedAt.Format(time.RFC3339Nano)
		}
		return domain.NewUsersResult(users, &cursor), err
	default:
		return nil, domain.ErrInvalidQuery
	}
}

func (s *AdminService) UpdateUser(ctx context.Context, token *domain.Token, update *domain.UserUpdate) error {
	if token.TokenType != domain.AccessToken {
		return domain.ErrInvalidTokenType
	}
	if token.UserRole != domain.Admin {
		return domain.ErrInvalidTokenRole
	}

	hasFieldToUpdate := false

	if update.Username != nil {
		hasFieldToUpdate = true
	}
	if update.Email != nil {
		hasFieldToUpdate = true
	}
	if update.Password != nil {
		hash, err := s.passwordHasher.Hash(*update.Password)
		if err != nil {
			return err
		}
		update.Password = &hash
		hasFieldToUpdate = true
	}
	if update.Role != nil {
		hasFieldToUpdate = true
	}

	if !hasFieldToUpdate {
		return domain.ErrNoFieldsToUpdate
	}

	if err := s.userRepository.UpdateUser(ctx, update); err != nil {
		return err
	}

	if err := s.tokenRepository.DeleteAllTokensByUserId(ctx, token.UserId); err != nil {
		return err
	}

	return nil
}
