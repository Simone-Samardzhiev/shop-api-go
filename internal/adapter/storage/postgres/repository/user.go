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
				"adding user failed",
				zap.String("id", user.Id.String()),
				zap.String("username", user.Username),
				zap.String("email", user.Email),
				zap.Error(err),
			)
		return domain.ErrInternal
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
					"fetching user failed",
					zap.String("username", username),
					zap.Error(err),
				)
			return nil, domain.ErrInternal
		}
	}
	return &user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at 
		FROM users
		WHERE id = $1`,
		id)

	var user domain.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		zap.L().
			Error(
				"fetching user failed",
				zap.String("id", id.String()),
				zap.Error(err),
			)
		return nil, domain.ErrInternal
	}
	return &user, nil
}

func (r *UserRepository) GetUsersByOffestPagination(ctx context.Context, page, limit int, role *domain.UserRole) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at 
		FROM users
		WHERE $1::user_role_enum IS NULL OR $1::user_role_enum = role
		OFFSET $2 LIMIT $3`,
		role,
		(page-1)*limit,
		limit,
	)
	if err != nil {
		zap.L().
			Error(
				"fetching users failed",
				zap.Int("page", page),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternal
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"closing rows failed",
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
					"error parsing row",
					zap.Error(scanErr),
				)
			return users, domain.ErrInternal
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUsersByTimePagination(ctx context.Context, after time.Time, limit int, role *domain.UserRole) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at FROM users
		WHERE created_at > $1 AND ($2::user_role_enum IS NULL OR $2::user_role_enum = role)
		ORDER BY created_at
		LIMIT $3`,
		after,
		role,
		limit,
	)

	if err != nil {
		zap.L().
			Error(
				"error fetching users",
				zap.Time("after", after),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternal
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"error closing rows",
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
					"error parsing row",
					zap.Error(err),
				)
			return nil, domain.ErrInternal
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SearchUserByUsername(ctx context.Context, username string, limit int, role *domain.UserRole) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at FROM users
		WHERE username % $1 AND ($2::user_role_enum IS NULL OR $2::user_role_enum = role)
		ORDER BY similarity(username, $1) DESC
		LIMIT $3`,
		username,
		role,
		limit,
	)

	if err != nil {
		zap.L().
			Error(
				"error searching for user",
				zap.String("username", username),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternal
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"error closing rows",
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
					"error parsing row",
					zap.Error(err),
				)
			return nil, domain.ErrInternal
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SearchUserByEmail(ctx context.Context, email string, limit int, role *domain.UserRole) ([]domain.User, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, email, role, created_at, updated_at FROM users
		WHERE similarity(email, $1) > 0.6 AND ($2::user_role_enum IS NULL OR $2::user_role_enum = role)
		ORDER BY similarity(email, $1) DESC
		limit $3`,
		email,
		role,
		limit,
	)

	if err != nil {
		zap.L().
			Error(
				"error searching for user",
				zap.String("email", email),
				zap.Int("limit", limit),
				zap.Error(err),
			)
		return nil, domain.ErrInternal
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			zap.L().
				Error(
					"error closing rows",
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
					"error parsing row",
					zap.Error(err),
				)
			return nil, domain.ErrInternal
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, update *domain.UserUpdate) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users 
		SET username = COALESCE($1, username),
		email = COALESCE($2, email),
		password = COALESCE($3, password),
		role = COALESCE($4, role),
		updated_at = now()
		WHERE id = $5`,
		update.Username,
		update.Email,
		update.Password,
		update.Role,
		update.Id,
	)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		switch {
		case pqErr.Constraint == "users_username_key":
			return domain.ErrUsernameAlreadyInUse
		case pqErr.Constraint == "users_email_key":
			return domain.ErrEmailAlreadyInUse
		}
	} else if err != nil {
		zapFields := make([]zap.Field, 0, 3)
		if update.Username != nil {
			zapFields = append(zapFields, zap.String("username", *update.Username))
		}
		if update.Email != nil {
			zapFields = append(zapFields, zap.String("email", *update.Email))
		}
		if update.Role != nil {
			zapFields = append(zapFields, zap.String("role", string(*update.Role)))
		}

		zap.L().
			Error(
				"error updating user",
				zapFields...,
			)
		return domain.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().
			Error(
				"error getting rows affected",
				zap.Error(err),
			)
		return domain.ErrInternal
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
