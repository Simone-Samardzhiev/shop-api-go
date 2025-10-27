package http

import (
	"encoding/base64"
	"net/http"
	"shop-api-go/internal/adapter/handler/http/request"
	"shop-api-go/internal/adapter/handler/http/response"
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

// GetUsers godoc
// @Summary      Users information
// @Description  Retrieves user information with optional filters and pagination. Requires admin privileges and a valid JWT token in the Authorization header.
// @Tags         Admin
// @Security     BearerAuth
// @Param        Authorization  header    string  true   "Bearer access token"
// @Param        id             query     string  false  "Filter by user ID (UUID)"
// @Param        username       query     string  false  "Filter by username"
// @Param        email          query     string  false  "Filter by email address"
// @Param        role           query     string  false  "Filter by user role (e.g. 'admin', 'user')"
// @Param        page           query     int     false  "Page number for pagination (min=1)"
// @Param        cursor         query     string  false  "Base64-encoded timestamp cursor for pagination"
// @Param        limit          query     int     false  "Maximum number of users to return (min=1)"
// @Produce      json
// @Success      200  {object}  response.FetchingUsersResponse "List of users"
// @Failure      400  {object}  response.ErrorResponse "Invalid query parameters"
// @Failure      401  {object}  response.ErrorResponse "Missing or invalid JWT token"
// @Failure      403  {object}  response.ErrorResponse "Insufficient permissions or invalid token type(expected access token)"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Router       /admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		response.HandleError(c, domain.ErrInternal)
		return
	}
	domainToken, ok := token.(*domain.Token)
	if !ok {
		response.HandleError(c, domain.ErrInternal)
		return
	}
	query := request.GetUserQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandleBindingError(c, err)
		return
	}

	var id *uuid.UUID
	if query.Id != nil {
		parsedId, err := uuid.Parse(*query.Id)
		if err != nil {
			response.HandleError(c, domain.ErrInvalidUUID)
			return
		}
		id = &parsedId
	}

	var after *time.Time
	if query.Cursor != nil {
		if *query.Cursor == "" {
			after = &time.Time{}
		} else {
			decoded, decodeErr := base64.URLEncoding.DecodeString(*query.Cursor)
			if decodeErr != nil {
				response.HandleError(c, domain.ErrInvalidCursor)
				return
			}
			parsedTime, decodeErr := time.Parse(time.RFC3339Nano, string(decoded))
			if decodeErr != nil {
				response.HandleError(c, domain.ErrInvalidCursor)
				return
			}
			after = &parsedTime
		}
	}

	result, err := h.adminService.GetUsers(
		c,
		domainToken,
		domain.NewGetUsers(
			id, query.Username, query.Email, query.Role, query.Page, after, query.Limit,
		),
	)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewFetchingUsersResponse(result))
}

// UpdateUser godoc
// @Summary      Update user by admin
// @Description  Allows an admin to update a user's information, including username, email, password, and role. Requires a valid admin JWT token.
// @Tags         Admin
// @Security     BearerAuth
// @Param        Authorization  header    string              true   "Bearer access token"
// @Param        id             path      string              true   "User ID (UUID) to update"
// @Param        request        body      request.UpdateUser  true   "Fields to update for the user"
// @Success      200            {string}  string              "User updated successfully"
// @Failure      400            {object}  response.ErrorResponse "Invalid request payload or parameters"
// @Failure      401            {object}  response.ErrorResponse "Unauthorized – invalid token"
// @Failure      403            {object}  response.ErrorResponse "Forbidden – insufficient permissions or invalid token type(expected access token)"
// @Failure      404            {object}  response.ErrorResponse "User not found"
// @Failure      500            {object}  response.ErrorResponse "Internal server error"
// @Router       /admin/users/update/{id} [patch]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		response.HandleError(c, domain.ErrInternal)
		return
	}
	domainToken, ok := token.(*domain.Token)
	if !ok {
		response.HandleError(c, domain.ErrInternal)
		return
	}

	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		response.HandleError(c, domain.ErrInvalidUUID)
		return
	}

	var req request.UpdateUser
	if err = c.ShouldBindJSON(&req); err != nil {
		response.HandleBindingError(c, err)
		return
	}

	if err = h.adminService.
		UpdateUser(
			c,
			domainToken,
			domain.NewUserUpdate(parsedId, req.Username, req.Email, req.Password, req.Role),
		); err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
