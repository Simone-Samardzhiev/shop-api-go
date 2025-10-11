package service

import (
	"context"
	"errors"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"

	"github.com/google/uuid"
)

// AuthService implements port.AuthService interface and provides access to admin-related business logic.
type AuthService struct {
	tokenGenerator  port.TokenGenerator
	passwordHasher  port.PasswordHasher
	tokenRepository port.TokenRepository
	userRepository  port.UserRepository
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(tokenGenerator port.TokenGenerator, passwordHasher port.PasswordHasher, tokenRepository port.TokenRepository, userRepository port.UserRepository) *AuthService {
	return &AuthService{
		tokenGenerator:  tokenGenerator,
		passwordHasher:  passwordHasher,
		tokenRepository: tokenRepository,
		userRepository:  userRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, user *domain.User) (*domain.TokenGroup, error) {
	fetchedUser, err := s.userRepository.GetUserByUsername(ctx, user.Username)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, domain.ErrWrongCredentials
	} else if err != nil {
		return nil, err
	}

	err = s.passwordHasher.Compare(user.Password, fetchedUser.Password)
	if err != nil {
		return nil, err
	}

	accessToken := domain.Token{
		Id:        uuid.New(),
		UserId:    fetchedUser.Id,
		TokenType: domain.AccessToken,
		UserRole:  fetchedUser.Role,
	}

	signedAccessToken, err := s.tokenGenerator.SignToken(&accessToken)
	if err != nil {
		return nil, err
	}

	refreshToken := domain.Token{
		Id:        uuid.New(),
		UserId:    fetchedUser.Id,
		TokenType: domain.RefreshToken,
		UserRole:  fetchedUser.Role,
	}
	signedRefreshToken, err := s.tokenGenerator.SignToken(&refreshToken)
	if err != nil {
		return nil, err
	}
	err = s.tokenRepository.AddToken(ctx, &refreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenGroup{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (s *AuthService) RefreshSession(ctx context.Context, token *domain.Token) (*domain.TokenGroup, error) {
	if token.TokenType != domain.RefreshToken {
		return nil, domain.ErrInvalidTokenType
	}

	err := s.tokenRepository.DeleteToken(ctx, token.Id)
	if errors.Is(err, domain.ErrTokenNotFound) {
		return nil, domain.ErrInvalidToken
	} else if err != nil {
		return nil, err
	}

	accessToken := domain.Token{
		Id:        uuid.New(),
		UserId:    token.UserId,
		TokenType: domain.AccessToken,
		UserRole:  token.UserRole,
	}
	signedAccessToken, err := s.tokenGenerator.SignToken(&accessToken)
	if err != nil {
		return nil, err
	}

	token = &domain.Token{
		Id:        uuid.New(),
		UserId:    token.UserId,
		TokenType: domain.RefreshToken,
		UserRole:  token.UserRole,
	}
	signedRefreshToken, err := s.tokenGenerator.SignToken(token)
	if err != nil {
		return nil, err
	}

	err = s.tokenRepository.AddToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &domain.TokenGroup{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}
