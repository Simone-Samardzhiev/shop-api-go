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

// NewUser create a new User instance.
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
