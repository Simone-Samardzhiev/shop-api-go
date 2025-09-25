package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shop-api-go/internal/core/domain"

	"github.com/lib/pq"
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

func (r *UserRepository) AddUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)`,
		user.Id,
		user.Username,
		user.Email,
		user.Password,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				switch pqErr.Constraint {
				case "users_username_key":
					return domain.ErrUsernameAlreadyInUse
				case "users_email_key":
					return domain.ErrEmailAlreadyInUse
				}
			}
		}

		return fmt.Errorf("%w: %v", domain.ErrInternalServerError, err)
	}

	return nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, password, role FROM users
                WHERE username = $1`,
		username,
	)

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrWrongCredentials
		} else {
			return nil, fmt.Errorf("%w: %v", domain.ErrInternalServerError, err)
		}
	}
	return &user, nil
}
