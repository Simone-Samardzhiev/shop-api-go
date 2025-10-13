package repository

import (
	"context"
	"database/sql"
	"errors"
	"shop-api-go/internal/core/domain"
	"time"

	"github.com/google/uuid"
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
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_key":
				return domain.ErrUsernameAlreadyInUse
			case "users_email_key":
				return domain.ErrEmailAlreadyInUse
			}
		}

		zap.L().
			Error(
				"postgres/UserRepository.AddUser failed",
				zap.String("id", user.Id.String()),
				zap.String("username", user.Username),
				zap.String("email", user.Email),
				zap.String("createdAt", user.CreatedAt.String()),
				zap.String("updatedAt", user.UpdatedAt.String()),
				zap.Error(err),
			)
		return domain.ErrInternalServerError
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
			return nil, domain.ErrUserNotFound
		} else {
			zap.L().
				Error(
					"postgres/UserRepository.GetUserByUsername failed",
					zap.String("username", username),
					zap.Error(err),
				)
			return nil, domain.ErrInternalServerError
		}
	}
	return &user, nil
}

func (r *UserRepository) GetUsersByOffestPagination(ctx context.Context, page, limit int) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at 
		FROM users
		OFFSET $1 LIMIT $2`,
		(page-1)*limit,
		limit,
	)
	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.GetUsersByOffestPagination failed",
				zap.Int("page", page),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternalServerError
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"postgres/UserRepository.GetUsersByOffestPagination failed closing rows",
					zap.Error(closeErr),
				)
		}
	}()

	users := make([]domain.User, 0, limit)
	for rows.Next() {
		var user domain.User
		scanErr := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if scanErr != nil {
			zap.L().
				Error(
					"postgres/UserRepository.GetUsersByOffestPagination failed scanning rwo",
					zap.Error(scanErr),
				)
			return users, domain.ErrInternalServerError
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUsersByTimePagination(ctx context.Context, after time.Time, limit int) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at FROM users
		WHERE created_at > $1
		ORDER BY created_at
		LIMIT $2`,
		after, limit,
	)

	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.GetUsersByTimePagination failed",
				zap.Time("after", after),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternalServerError
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"postgres/UserRepository.GetUsersByTimePagination failed closing rows",
					zap.Error(closeErr),
				)
		}
	}()

	users := make([]domain.User, 0, limit)
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			zap.L().
				Error(
					"postgres/UserRepository.GetUsersByTimePagination failed scanning row",
					zap.Error(err),
				)
			return nil, domain.ErrInternalServerError
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SearchUserByUsername(ctx context.Context, username string, limit int) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at FROM users
		WHERE username % $1
		ORDER BY similarity(username, $1) DESC
		LIMIT $2`,
		username,
		limit,
	)

	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.SearchUserByUsername failed",
				zap.String("username", username),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternalServerError
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"postgres/UserRepository.SearchUserByUsername failed closing rows",
					zap.Error(closeErr),
				)
		}
	}()

	users := make([]domain.User, 0, limit)
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			zap.L().
				Error(
					"postgres/UserRepository.SearchUserByUsername failed scanning row",
					zap.Error(err),
				)
			return nil, domain.ErrInternalServerError
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SearchUserByEmail(ctx context.Context, email string, limit int) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at FROM users
		WHERE similarity(email, $1) > 0.6
		ORDER BY similarity(email, $1) DESC
		limit $2`,
		email,
		limit,
	)

	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.SearchUserByEmail failed",
				zap.String("email", email),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternalServerError
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"postgres/UserRepository.SearchUserByEmail failed closing rows",
					zap.Error(closeErr),
				)
		}
	}()

	users := make([]domain.User, 0, limit)
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			zap.L().
				Error(
					"postgres/UserRepository.SearchUserByEmail failed scanning row",
					zap.Error(err),
				)
			return nil, domain.ErrInternalServerError
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at 
		FROM users WHERE id = $1`,
		id)

	var user domain.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.GetUserById failed",
				zap.String("id", id.String()),
				zap.Error(err),
			)
		return nil, domain.ErrInternalServerError
	}
	return &user, nil
}

func (r *UserRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users 
		SET username = $1, updated_at = now()
		WHERE id = $2`,
		username,
		id,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			return domain.ErrUsernameAlreadyInUse
		}

		zap.L().
			Error(
				"postgres/UserRepository.UpdateUsername failed",
				zap.String("id", id.String()),
				zap.String("username", username),
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.UpdateUsername failed getting rows affected",
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users
		SET email = $1, updated_at = now() 
		WHERE id = $2`,
		email,
		id,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			return domain.ErrEmailAlreadyInUse
		}

		zap.L().
			Error(
				"postgres/UserRepository.UpdateEmail failed",
				zap.String("id", id.String()),
				zap.String("email", email),
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.UpdateEmail failed getting affected rows",
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users
		SET password = $1, updated_at = now() 
		WHERE id = $2`,
		password,
		id,
	)
	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.UpdatePassword failed",
				zap.String("id", id.String()),
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.UpdatePassword getting rows affected",
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role domain.UserRole) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET role = $1, updated_at = now()
		WHERE id = $2`,
		role,
		id,
	)

	if err != nil {
		zap.L().Error(
			"postgres/UserRepository.UpdateRole failed",
			zap.String("id", id.String()),
			zap.String("role", string(role)),
			zap.Error(err),
		)
		return domain.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().
			Error(
				"postgres/UserRepository.UpdateRole failed getting rows affected",
				zap.Error(err),
			)
		return domain.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
