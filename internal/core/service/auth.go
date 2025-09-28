package service

import (
	"context"
	"errors"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"shop-api-go/internal/core/util"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
	if errors.Is(err, domain.ErrWrongCredentials) {
		return nil, domain.ErrWrongCredentials
	} else if err != nil {
		zap.L().Error("AuthService.Login failed",
			zap.String("username", user.Username),
			zap.Error(err))
		return nil, domain.ErrInternalServerError
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
		zap.L().Error("AuthService.Login failed",
			zap.String("username", user.Username),
			zap.Error(err),
		)

		return nil, domain.ErrInternalServerError
	}

	refreshToken := domain.Token{
		Id:        uuid.New(),
		UserId:    fetchedUser.Id,
		TokenType: domain.RefreshToken,
		UserRole:  fetchedUser.Role,
	}
	signedRefreshToken, err := as.tokenGenerator.SignToken(&refreshToken)
	if err != nil {
		zap.L().Error("AuthService.Login failed",
			zap.String("username", user.Username),
			zap.Error(err),
		)
	}
	err = as.tokenRepository.AddToken(ctx, &refreshToken)
	if err != nil {
		zap.L().Error("AuthService.Login failed",
			zap.String("username", user.Username),
			zap.Error(err),
		)
		return nil, domain.ErrInternalServerError
	}

	return &domain.TokenGroup{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}
