package repository

import (
	"context"
	"database/sql"
	"shop-api-go/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TokenRepository implements port.TokenRepository and provides
// access to postgres database
type TokenRepository struct {
	db *sql.DB
}

// NewTokenRepository creates new TokenRepository instance.
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (t *TokenRepository) AddToken(ctx context.Context, token *domain.Token) error {
	_, err := t.db.ExecContext(
		ctx,
		`INSERT INTO tokens(id, user_id, token_type, expires)
		VALUES ($1, $2, $3, $4)`,
		token.Id,
		token.UserId,
		token.TokenType,
		token.ExpiresAt,
	)
	if err != nil {
		zap.L().
			Error(
				"postgres/TokenRepository.AddToken failed",
				zap.String("tokenId", token.Id.String()),
				zap.String("userId", token.UserId.String()),
				zap.String("tokenType", string(token.TokenType)),
				zap.String("userRole", string(token.UserRole)),
				zap.String("expiresAt", token.ExpiresAt.String()),
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}
	return nil
}

func (t *TokenRepository) DeleteToken(ctx context.Context, id uuid.UUID) error {
	result, err := t.db.ExecContext(ctx, "DELETE FROM tokens WHERE id = $1", id)
	if err != nil {
		zap.L().
			Error(
				"postgres/TokenRepository.DeleteToken failed",
				zap.String("id", id.String()),
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().
			Error(
				"postgres/TokenRepository.DeleteToken failed getting affected rows",
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return domain.ErrTokenNotFound
	}
	return nil
}

func (t *TokenRepository) DeleteAllTokensByUserId(ctx context.Context, userId uuid.UUID) error {
	_, err := t.db.ExecContext(ctx, "DELETE FROM tokens WHERE user_id = $1", userId)
	if err != nil {
		zap.L().
			Error(
				"postgres/TokenRepository.DeleteAllTokensByUserId failed",
				zap.String("userId", userId.String()),
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}
	return nil
}

func (t *TokenRepository) DeleteExpiredTokens() error {
	_, err := t.db.Exec("DELETE FROM tokens WHERE expires < NOW()")
	if err != nil {
		zap.L().
			Error(
				"postgres/TokenRepository.DeleteExpiredTokens failed",
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}
	return nil
}
