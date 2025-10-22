package http

import (
	"net/http"
	"shop-api-go/internal/adapter/handler/http/request"
	"shop-api-go/internal/adapter/handler/http/response"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/service"

	"github.com/gin-gonic/gin"
)

// UserHandler represents HTTP handler for user-related requests.
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler returns a new UserHandler instance.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleBindingError(c, err)
		return
	}

	err := h.userService.Register(c, &domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *UserHandler) UpdateAccount(c *gin.Context) {
	var req request.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleBindingError(c, err)
		return
	}

	if err := h.userService.
		UpdateAccount(
			c,
			domain.NewUpdateAccount(
				req.Username,
				req.Password,
				req.NewUsername,
				req.NewEmail,
				req.NewPassword),
		); err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
