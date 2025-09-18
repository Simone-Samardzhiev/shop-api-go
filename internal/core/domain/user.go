package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole is an emum for user's role.
type UserRole string

// UserRole enum values.
const (
	Admin     UserRole = "Admin"
	Client    UserRole = "Client"
	Delivery  UserRole = "Delivery"
	Warehouse UserRole = "Warehouse"
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
