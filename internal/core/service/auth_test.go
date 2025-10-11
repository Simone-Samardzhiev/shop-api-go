package service_test

import (
	"context"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port/mock"
	"shop-api-go/internal/core/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name               string
		user               *domain.User
		expectedError      error
		expectedTokenGroup *domain.TokenGroup
		mockSetup          func(
			mockTokenGenerator *mock.MockTokenGenerator,
			mockPasswordHasher *mock.MockPasswordHasher,
			mockTokenRepository *mock.MockTokenRepository,
			mockUserRepository *mock.MockUserRepository,
		)
	}{
		{
			name: "success",
			user: &domain.User{
				Password: "password",
			},
			expectedTokenGroup: &domain.TokenGroup{
				AccessToken:  "token",
				RefreshToken: "token",
			},
			expectedError: nil,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {

				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(""),
						).
						Return(&domain.User{
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.EXPECT().
						Compare("password", "hashedPassword").
						Return(nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("token", nil).
						Times(2),
					mockTokenRepository.
						EXPECT().
						AddToken(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(&domain.Token{}),
						).
						Return(nil),
				)
			},
		}, {
			name:               "user not found",
			user:               &domain.User{},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrWrongCredentials,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(""),
					).
					Return(nil, domain.ErrUserNotFound)
			},
		}, {
			name:               "error fetching user",
			user:               &domain.User{},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository,
			) {
				mockUserRepository.
					EXPECT().
					GetUserByUsername(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf("")).
					Return(nil, domain.ErrInternalServerError)
			},
		}, {
			name: "password mismatch",
			user: &domain.User{
				Password: "password",
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrWrongCredentials,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(""),
						).Return(&domain.User{
						Password: "hashedPassword",
					}, nil),
					mockPasswordHasher.EXPECT().
						Compare("password", "hashedPassword").
						Return(domain.ErrWrongCredentials),
				)
			},
		}, {
			name: "error signing access token",
			user: &domain.User{
				Password: "password",
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(""),
						).
						Return(&domain.User{
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.EXPECT().
						Compare("password", "hashedPassword").
						Return(nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("", domain.ErrInternalServerError),
				)
			},
		}, {
			name: "error signing refresh token",
			user: &domain.User{
				Password: "password",
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(""),
						).
						Return(&domain.User{
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.EXPECT().
						Compare("password", "hashedPassword").
						Return(nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("token", nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("", domain.ErrInternalServerError),
				)
			},
		}, {
			name: "error adding token",
			user: &domain.User{
				Password: "password",
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				gomock.InOrder(
					mockUserRepository.
						EXPECT().
						GetUserByUsername(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(""),
						).
						Return(&domain.User{
							Password: "hashedPassword",
						}, nil),
					mockPasswordHasher.EXPECT().
						Compare("password", "hashedPassword").
						Return(nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("token", nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("token", nil),

					mockTokenRepository.
						EXPECT().
						AddToken(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(&domain.Token{}),
						).
						Return(domain.ErrInternalServerError),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTokenGenerator := mock.NewMockTokenGenerator(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			tt.mockSetup(mockTokenGenerator, mockPasswordHasher, mockTokenRepository, mockUserRepository)

			tokenGroup, err := service.
				NewAuthService(mockTokenGenerator, mockPasswordHasher, mockTokenRepository, mockUserRepository).
				Login(context.Background(), tt.user)

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedTokenGroup, tokenGroup)
		})
	}
}

func TestAuthService_RefreshSession(t *testing.T) {
	tests := []struct {
		name               string
		token              *domain.Token
		expectedTokenGroup *domain.TokenGroup
		expectedError      error
		mockSetup          func(
			mockTokenGenerator *mock.MockTokenGenerator,
			mockPasswordHasher *mock.MockPasswordHasher,
			mockTokenRepository *mock.MockTokenRepository,
			mockUserRepository *mock.MockUserRepository,
		)
	}{
		{
			name: "success",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedTokenGroup: &domain.TokenGroup{
				AccessToken:  "token",
				RefreshToken: "token",
			},
			expectedError: nil,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				mockTokenRepository.
					EXPECT().
					DeleteToken(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
					).Return(nil)
				mockTokenGenerator.
					EXPECT().
					SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
					Return("token", nil).Times(2)
				mockTokenRepository.
					EXPECT().
					AddToken(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(&domain.Token{})).
					Return(nil)
			},
		}, {
			name: "wrong token type",
			token: &domain.Token{
				TokenType: domain.AccessToken,
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInvalidTokenType,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository,
			) {

			},
		}, {
			name: "error deleting token",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				mockTokenRepository.
					EXPECT().
					DeleteToken(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
					).
					Return(domain.ErrInternalServerError)
			},
		}, {
			name: "error signing access token",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				gomock.InOrder(
					mockTokenRepository.
						EXPECT().
						DeleteToken(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{})).
						Return(nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("", domain.ErrInternalServerError),
				)
			},
		}, {
			name: "failed to signing refresh token",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				gomock.InOrder(
					mockTokenRepository.
						EXPECT().
						DeleteToken(
							gomock.AssignableToTypeOf(context.Background()),
							gomock.AssignableToTypeOf(uuid.UUID{})).
						Return(nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("token", nil),
					mockTokenGenerator.
						EXPECT().
						SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
						Return("", domain.ErrInternalServerError),
				)
			},
		}, {
			name: "error adding token",
			token: &domain.Token{
				TokenType: domain.RefreshToken,
			},
			expectedTokenGroup: nil,
			expectedError:      domain.ErrInternalServerError,
			mockSetup: func(
				mockTokenGenerator *mock.MockTokenGenerator,
				mockPasswordHasher *mock.MockPasswordHasher,
				mockTokenRepository *mock.MockTokenRepository,
				mockUserRepository *mock.MockUserRepository) {
				mockTokenRepository.
					EXPECT().
					DeleteToken(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(uuid.UUID{}),
					).Return(nil)
				mockTokenGenerator.
					EXPECT().
					SignToken(gomock.AssignableToTypeOf(&domain.Token{})).
					Return("token", nil).Times(2)
				mockTokenRepository.
					EXPECT().
					AddToken(
						gomock.AssignableToTypeOf(context.Background()),
						gomock.AssignableToTypeOf(&domain.Token{})).
					Return(domain.ErrInternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTokenGenerator := mock.NewMockTokenGenerator(ctrl)
			mockPasswordHasher := mock.NewMockPasswordHasher(ctrl)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			tt.mockSetup(mockTokenGenerator, mockPasswordHasher, mockTokenRepository, mockUserRepository)

			tokenGroup, err := service.NewAuthService(
				mockTokenGenerator,
				mockPasswordHasher,
				mockTokenRepository,
				mockUserRepository,
			).
				RefreshSession(context.Background(), tt.token)

			if tt.expectedError != nil {
				require.ErrorIs(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedTokenGroup, tokenGroup)
		})
	}
}
