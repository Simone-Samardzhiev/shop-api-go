package repository

import (
	"context"
	"database/sql"
	"shop-api-go/internal/core/domain"
)

// UserRepository implements port.UserRepository and provides
// access to postgres database.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)`,
		user.Id,
		user.Username,
		user.Email,
		user.Password,
	)

	return err
}
