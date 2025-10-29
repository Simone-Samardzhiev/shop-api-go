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
				"failed to insert token",
				zap.String("tokenId", token.Id.String()),
				zap.String("userId", token.UserId.String()),
				zap.String("tokenType", string(token.TokenType)),
				zap.String("userRole", string(token.UserRole)),
				zap.Error(err),
			)
		return domain.ErrInternal
	}
	return nil
}

func (t *TokenRepository) DeleteToken(ctx context.Context, id uuid.UUID) error {
	result, err := t.db.ExecContext(ctx, "DELETE FROM tokens WHERE id = $1", id)
	if err != nil {
		zap.L().
			Error(
				"failed to delete token",
				zap.String("id", id.String()),
				zap.Error(err),
			)
		return domain.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().
			Error(
				"failed to retrieve rows affected",
				zap.Error(err),
			)
		return domain.ErrInternal
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
				"failed to delete all tokens",
				zap.String("userId", userId.String()),
				zap.Error(err),
			)
		return domain.ErrInternal
	}
	return nil
}

func (t *TokenRepository) DeleteExpiredTokens() error {
	_, err := t.db.Exec("DELETE FROM tokens WHERE expires < NOW()")
	if err != nil {
		zap.L().
			Error(
				"failed to delete expired tokens",
				zap.Error(err),
			)
		return domain.ErrInternal
	}
	return nil
}
