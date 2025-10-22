package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole is an emum for user's role.
type UserRole string

// UserRole enum values.
const (
	Admin     UserRole = "admin"
	Client    UserRole = "client"
	Delivery  UserRole = "delivery"
	Warehouse UserRole = "warehouse"
)

// User is an entity representing a user.
type User struct {
	Id        uuid.UUID
	Username  string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User instance.
func NewUser(id uuid.UUID, username, email, password string, role UserRole, createdAt, UpdatedAt time.Time) *User {
	return &User{
		Id:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: UpdatedAt,
	}
}

// UserUpdate is an DTO for updating user's fields.
type UserUpdate struct {
	Id       uuid.UUID
	Username *string
	Email    *string
	Password *string
	Role     *UserRole
}

// NewUserUpdate creates a new UserUpdate instance.
func NewUserUpdate(id uuid.UUID, username, email, password *string, role *UserRole) *UserUpdate {
	return &UserUpdate{
		Id:       id,
		Username: username,
		Email:    email,
		Password: password,
		Role:     role,
	}
}

// GetUsers is a DTO for getting users info.
type GetUsers struct {
	Id       *uuid.UUID
	Username *string
	Email    *string
	Role     *UserRole
	Page     *int
	After    *time.Time
	Limit    *int
}

// NewGetUsers creates a new GetUsers instance.
func NewGetUsers(id *uuid.UUID, username, email *string, role *UserRole, page *int, after *time.Time, limit *int) *GetUsers {
	return &GetUsers{
		Id:       id,
		Username: username,
		Email:    email,
		Role:     role,
		Page:     page,
		After:    after,
		Limit:    limit,
	}
}

// UpdateAccount is a DTO for a user to update his account.
type UpdateAccount struct {
	Username    string
	Password    string
	NewUsername *string
	NewEmail    *string
	NewPassword *string
}

// NewUpdateAccount creates new UpdateAccount instance.
func NewUpdateAccount(username, password string, newUsername, newEmail, newPassword *string) *UpdateAccount {
	return &UpdateAccount{
		Username:    username,
		Password:    password,
		NewUsername: newUsername,
		NewEmail:    newEmail,
		NewPassword: newPassword,
	}
}

// UsersResult is a DTO for fetching users result.
type UsersResult struct {
	Users  []User
	Cursor *string
}

// NewUsersResult creates a new UsersResult instance.
func NewUsersResult(users []User, cursor *string) *UsersResult {
	return &UsersResult{
		Users:  users,
		Cursor: cursor,
	}
}
