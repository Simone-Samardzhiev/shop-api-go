package http

import (
	"encoding/base64"
	"net/http"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminHandler represent HTTP handler for admin-related requests.
type AdminHandler struct {
	adminService port.AdminService
}

// NewAdminHandler creates a new AdminHandler instance.
func NewAdminHandler(userService port.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: userService,
	}
}

type userResponse struct {
	ID        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type getUserResponse struct {
	Users  []userResponse `json:"users"`
	Cursor *string        `json:"cursor"`
}

func newGetUsersResponse(result *domain.UsersResult) *getUserResponse {
	users := make([]userResponse, 0, len(result.Users))
	for _, user := range result.Users {
		users = append(users, userResponse{
			ID:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	if result.Cursor != nil {
		encodedCursor := base64.URLEncoding.EncodeToString([]byte(*result.Cursor))
		result.Cursor = &encodedCursor
	}

	return &getUserResponse{
		Users:  users,
		Cursor: result.Cursor,
	}
}

type getUsersQueryParams struct {
	Id       *uuid.UUID       `form:"id"`
	Username *string          `form:"username"`
	Email    *string          `form:"email"`
	Role     *domain.UserRole `form:"role" binding:"omitempty,user_role"`
	Page     *int             `form:"page" binding:"omitempty,min=1"`
	Cursor   *string          `form:"cursor"`
	Limit    *int             `form:"limit" binding:"omitempty,min=1"`
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		handleError(c, domain.ErrInternalServerError)
		return
	}

	domainToken, ok := token.(*domain.Token)
	if !ok {
		handleError(c, domain.ErrInternalServerError)
		return
	}

	query := getUsersQueryParams{}
	if err := c.BindQuery(&query); err != nil {
		handleBindingError(c, err)
		return
	}

	var after *time.Time
	if query.Cursor != nil {
		if *query.Cursor == "" {
			after = &time.Time{}
		} else {
			decoded, err := base64.URLEncoding.DecodeString(*query.Cursor)
			if err != nil {
				handleError(c, domain.ErrInvalidCursorFormat)
				return
			}
			parsedTime, err := time.Parse(time.RFC3339Nano, string(decoded))
			if err != nil {
				handleError(c, domain.ErrInvalidCursorFormat)
				return
			}
			after = &parsedTime
		}
	}

	response, err := h.adminService.GetUsers(
		c,
		domainToken,
		domain.NewGetUsers(
			query.Id, query.Username, query.Email, query.Role, query.Page, after, query.Limit,
		),
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, newGetUsersResponse(response))
}

type updateUserRequest struct {
	Id       uuid.UUID        `json:"id" binding:"required"`
	Username *string          `json:"username" binding:"omitempty,min_bytes=8,max_bytes=255"`
	Email    *string          `json:"email" binding:"omitempty,min_bytes=8,max_bytes=255"`
	Password *string          `json:"password" binding:"omitempty,password"`
	Role     *domain.UserRole `json:"role" binding:"omitempty,user_role"`
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		handleError(c, domain.ErrInternalServerError)
		return
	}

	domainToken, ok := token.(*domain.Token)
	if !ok {
		handleError(c, domain.ErrInternalServerError)
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	if err := h.adminService.
		UpdateUser(
			c,
			domainToken,
			domain.NewUserUpdate(req.Id, req.Username, req.Email, req.Password, req.Role),
		); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
