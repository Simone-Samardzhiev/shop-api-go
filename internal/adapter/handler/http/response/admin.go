package response

import (
	"encoding/base64"
	"shop-api-go/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

// user represents a response with user's information.
type user struct {
	Id        uuid.UUID       `json:"id" example:"1bd70616-480b-47b9-91f5-292b4f4a45b1"`
	Username  string          `json:"username" example:"Viktor123"`
	Email     string          `json:"email" example:"viktor.stavchev@gmail.com"`
	Role      domain.UserRole `json:"role" example:"client"`
	CreatedAt time.Time       `json:"createdAt" example:"2025-10-15T12:37:42.664482Z"`
	UpdatedAt time.Time       `json:"updatedAt" example:"2025-10-15T12:37:42.664482Z"`
}

// newUser creates a new user instance.
func newUser(u *domain.User) user {
	return user{
		Id:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FetchingUsersResponse represents a response when fetching users.
type FetchingUsersResponse struct {
	Users  []user  `json:"users"`
	Cursor *string `json:"cursor" example:"MjAyNS0xMC0xNVQxMjo0MDoxOS41NTU4Mjda"`
}

// NewFetchingUsersResponse creates a new FetchingUsersResponse instance.
func NewFetchingUsersResponse(result *domain.UsersResult) FetchingUsersResponse {
	users := make([]user, 0, len(result.Users))
	for _, u := range result.Users {
		users = append(users, newUser(&u))
	}

	if result.Cursor != nil {
		encodedCursor := base64.URLEncoding.EncodeToString([]byte(*result.Cursor))
		result.Cursor = &encodedCursor
	}
	return FetchingUsersResponse{
		Users:  users,
		Cursor: result.Cursor,
	}
}
