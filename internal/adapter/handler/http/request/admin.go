package request

import (
	"shop-api-go/internal/core/domain"
)

// GetUserQuery represents query parameters for fetching user.
type GetUserQuery struct {
	Id       *string          `form:"id"`
	Username *string          `form:"username"`
	Email    *string          `form:"email"`
	Role     *domain.UserRole `form:"role" binding:"omitempty,user_role"`
	Page     *int             `form:"page" binding:"omitempty,min=1"`
	Cursor   *string          `form:"cursor"`
	Limit    *int             `form:"limit" binding:"omitempty,min=1"`
}

// UpdateUser represents update user request body.
type UpdateUser struct {
	Username *string          `json:"username" binding:"omitempty,min_bytes=8,max_bytes=255"`
	Email    *string          `json:"email" binding:"omitempty,min_bytes=8,max_bytes=255"`
	Password *string          `json:"password" binding:"omitempty,password"`
	Role     *domain.UserRole `json:"role" binding:"omitempty,user_role"`
}
