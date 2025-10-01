package http

import (
	"net/http"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler represents HTTP handler for user-related requests.
type UserHandler struct {
	us *service.UserService
}

// NewUserHandler returns a new UserHandler instance.
func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		us: us,
	}
}

// registerRequest represent a request body for creating a user.
type registerRequest struct {
	Email    string `json:"email" binding:"required,email,min_bytes=8,max_bytes=255"`
	Username string `json:"username" binding:"required,min_bytes=8,max_bytes=255"`
	Password string `json:"password" binding:"required,password"`
}

func (uh *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	err := uh.us.Register(c, &domain.User{
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

// userInfoResponse contains users data that is safe to be shared.
type userInfoResponse struct {
	Id        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// mapUsersToUsersInfoResponse transforms a slice of domain.User to a slice of userInfoResponse.
func mapUsersToUsersInfoResponse(users []domain.User) []userInfoResponse {
	usersResponse := make([]userInfoResponse, len(users))
	for i, user := range users {
		usersResponse[i] = userInfoResponse{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return usersResponse
}

// getUsersByOffestPagination contains metadata used for user offest pagination.
type getUsersByOffestPagination struct {
	Limit int `json:"limit" binding:"min=1"`
	Page  int `json:"page" binding:"min=1"`
}

func (uh *UserHandler) GetUsersByPages(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		handleError(c, domain.ErrInternalServerError)
		return
	}

	domainToken, ok := token.(*domain.Token)
	if !ok {
		handleError(c, domain.ErrInternalServerError)
	}

	var req getUsersByOffestPagination
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	users, err := uh.us.GetUsersByOffestPagination(c, domainToken, req.Page, req.Limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapUsersToUsersInfoResponse(users))
}
