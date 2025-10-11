package service

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"shop-api-go/internal/core/util"

	"github.com/google/uuid"
)

// AuthService implements port.AuthService interface and provides access to admin-related business logic.
type AuthService struct {
	tokenGenerator  port.TokenGenerator
	tokenRepository port.TokenRepository
	userRepository  port.UserRepository
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(tokenGenerator port.TokenGenerator, tokenRepository port.TokenRepository, userRepository port.UserRepository) *AuthService {
	return &AuthService{
		tokenGenerator:  tokenGenerator,
		tokenRepository: tokenRepository,
		userRepository:  userRepository,
	}
}

func (as *AuthService) Login(ctx context.Context, user *domain.User) (*domain.TokenGroup, error) {
	fetchedUser, err := as.userRepository.GetUserByUsername(ctx, user.Username)
	fetchedUser, err := s.userRepository.GetUserByUsername(ctx, user.Username)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, domain.ErrWrongCredentials
	} else if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	if !util.ComparePassword(user.Password, fetchedUser.Password) {
		return nil, domain.ErrWrongCredentials
	}

	accessToken := domain.Token{
		Id:        uuid.New(),
		UserId:    fetchedUser.Id,
		TokenType: domain.AccessToken,
		UserRole:  fetchedUser.Role,
	}

	signedAccessToken, err := as.tokenGenerator.SignToken(&accessToken)
	if err != nil {
		return nil, err
	}

	refreshToken := domain.Token{
		Id:        uuid.New(),
		UserId:    fetchedUser.Id,
		TokenType: domain.RefreshToken,
		UserRole:  fetchedUser.Role,
	}
	signedRefreshToken, err := as.tokenGenerator.SignToken(&refreshToken)
	if err != nil {
		return nil, err
	}
	err = as.tokenRepository.AddToken(ctx, &refreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenGroup{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (as *AuthService) RefreshSession(ctx context.Context, token *domain.Token) (*domain.TokenGroup, error) {
	if token.TokenType != domain.RefreshToken {
		return nil, domain.ErrInvalidTokenType
	}

	result, err := as.tokenRepository.DeleteToken(ctx, token.Id)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, domain.ErrInvalidToken
	}

	accessToken := domain.Token{
		Id:        uuid.New(),
		UserId:    token.UserId,
		TokenType: domain.AccessToken,
		UserRole:  token.UserRole,
	}
	signedAccessToken, err := as.tokenGenerator.SignToken(&accessToken)
	if err != nil {
		return nil, err
	}

	token = &domain.Token{
		Id:        uuid.New(),
		UserId:    token.UserId,
		TokenType: domain.RefreshToken,
		UserRole:  token.UserRole,
	}
	signedRefreshToken, err := as.tokenGenerator.SignToken(token)
	if err != nil {
		return nil, err
	}

	err = as.tokenRepository.AddToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &domain.TokenGroup{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}
