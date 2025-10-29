package http

import (
	"net/http"
	"shop-api-go/internal/adapter/handler/http/request"
	"shop-api-go/internal/adapter/handler/http/response"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to user operations.
type UserHandler struct {
	userService port.UserService
}

// NewUserHandler returns a new UserHandler instance.
func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Registers a new user in the system using email, username, and password. Returns HTTP 201 on success.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      request.RegisterRequest  true  "Registration details"
// @Success      201                                               "User created successfully"
// @Failure      400      {object}  response.ErrorResponse "Invalid request payload or missing fields"
// @Failure      409      {object}  response.ErrorResponse "Email or username already exists"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Router       /users/register [post]
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

// UpdateAccount godoc
// @Summary      Update user account
// @Description  Updates a user's account. Requires current username and password for authentication. Optional fields include new username, new email, and new password.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request        body      request.UpdateAccountRequest true  "Update account payload"
// @Success      200                                                   "Account updated successfully"
// @Failure      400            {object}  response.ErrorResponse       "Invalid request payload or missing fields"
// @Failure      401            {object}  response.ErrorResponse       "Unauthorized - invalid credentials"
// @Failure      500            {object}  response.ErrorResponse       "Internal server error"
// @Router       /users/me/update [patch]
func (h *UserHandler) UpdateAccount(c *gin.Context) {
	var req request.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleBindingError(c, err)
		return
	}

	if err := h.userService.UpdateAccount(
		c,
		domain.NewUpdateAccount(
			req.Username,
			req.Password,
			req.NewUsername,
			req.NewEmail,
			req.NewPassword,
		),
	); err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
