package repository

import (
	"context"
	"database/sql"
	"errors"
	"shop-api-go/internal/core/domain"

	"github.com/lib/pq"
	"go.uber.org/zap"
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
				default:
					zap.L().Error("postgres/UserRepository.AddUser failed",
						zap.String("id", user.Id.String()),
						zap.String("username", user.Username),
						zap.String("email", user.Email),
						zap.String("createdAt", user.CreatedAt.String()),
						zap.String("updatedAt", user.UpdatedAt.String()),
						zap.Error(err),
					)
					return domain.ErrInternalServerError
				}
			}
		} else {
			zap.L().Error("postgres/UserRepository.AddUser failed", zap.Error(err))
			return domain.ErrInternalServerError
		}
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
			zap.L().Error("postgres/UserRepository.GetUserByUsername failed",
				zap.String("username", username),
				zap.Error(err),
			)
			return nil, domain.ErrInternalServerError
		}
	}
	return &user, nil
}
