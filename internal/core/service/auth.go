package service

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"shop-api-go/internal/core/util"

	"github.com/google/uuid"
)

type AuthService struct {
	tokenGenerator  port.TokenGenerator
	tokenRepository port.TokenRepository
	userRepository  port.UserRepository
}

func NewAuthService(tokenGenerator port.TokenGenerator, tokenRepository port.TokenRepository, userRepository port.UserRepository) *AuthService {
	return &AuthService{
		tokenGenerator:  tokenGenerator,
		tokenRepository: tokenRepository,
		userRepository:  userRepository,
	}
}

func (as *AuthService) Login(ctx context.Context, user *domain.User) (*domain.TokenGroup, error) {
	fetchedUser, err := as.userRepository.GetUserByUsername(ctx, user.Username)
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
