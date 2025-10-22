package response

import (
	"encoding/base64"
	"shop-api-go/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

// user represents a response with user's information.
type user struct {
	Id        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
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
	Cursor *string `json:"cursor"`
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
