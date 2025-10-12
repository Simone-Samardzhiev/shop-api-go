package http

import (
	"net/http"
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

// registerRequest represent a request body for creating a user.
type registerRequest struct {
	Email    string `json:"email" binding:"required,email,min_bytes=8,max_bytes=255"`
	Username string `json:"username" binding:"required,min_bytes=8,max_bytes=255"`
	Password string `json:"password" binding:"required,password"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	err := h.userService.Register(c, &domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// changeUsernameRequest represent a request body for changing username.
type changeUsernameRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewUsername string `json:"newUsername" binding:"required,min_bytes=8,max_bytes=255"`
}

func (h *UserHandler) ChangeUsername(c *gin.Context) {
	var req changeUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	err := h.userService.ChangeUsername(
		c,
		&domain.User{
			Username: req.Username,
			Password: req.Password,
		},
		req.NewUsername,
	)
	if err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusCreated)
}
