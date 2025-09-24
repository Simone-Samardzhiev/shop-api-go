package repository

import (
	"context"
	"database/sql"
	"fmt"
	"shop-api-go/internal/core/domain"
)

// TokenRepository implements port.TokenRepository and provides
// access to postgres database
type TokenRepository struct {
	db *sql.DB
}

// NewTaskRepository creates new TokenRepository instance.
func NewTaskRepository(db *sql.DB) *TokenRepository {
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
		return fmt.Errorf("%w: %v", domain.ErrInternalServerError, err)
	}
	return nil
}

func (t *TokenRepository) DeleteExpiredToken() error {
	_, err := t.db.Exec("DELETE FROM tokens WHERE expires < NOW()")
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInternalServerError, err)
	}
	return nil
}
