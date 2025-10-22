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

func (h *AdminHandler) GetUsers(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		response.HandleError(c, domain.ErrInternalServerError)
		return
	}
	domainToken, ok := token.(*domain.Token)
	if !ok {
		response.HandleError(c, domain.ErrInternalServerError)
		return
	}
	query := request.GetUserQuery{}
	if err := c.BindQuery(&query); err != nil {
		response.HandleBindingError(c, err)
		return
	}

	var after *time.Time
	if query.Cursor != nil {
		if *query.Cursor == "" {
			after = &time.Time{}
		} else {
			decoded, err := base64.URLEncoding.DecodeString(*query.Cursor)
			if err != nil {
				response.HandleError(c, domain.ErrInvalidCursorFormat)
				return
			}
			parsedTime, err := time.Parse(time.RFC3339Nano, string(decoded))
			if err != nil {
				response.HandleError(c, domain.ErrInvalidCursorFormat)
				return
			}
			after = &parsedTime
		}
	}

	result, err := h.adminService.GetUsers(
		c,
		domainToken,
		domain.NewGetUsers(
			query.Id, query.Username, query.Email, query.Role, query.Page, after, query.Limit,
		),
	)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewFetchingUsersResponse(result))
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		response.HandleError(c, domain.ErrInternalServerError)
		return
	}
	domainToken, ok := token.(*domain.Token)
	if !ok {
		response.HandleError(c, domain.ErrInternalServerError)
		return
	}

	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		response.HandleError(c, domain.ErrInvalidParam)
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
