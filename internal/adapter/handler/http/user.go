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

// userInfoResponse contains users data that is safe to be shared.
type userInfoResponse struct {
	Id        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// mapUsersToUserInfoResponse transforms a slice of domain.User to a slice of userInfoResponse.
func mapUsersToUserInfoResponse(users []domain.User) []userInfoResponse {
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

// usersByOffestPaginationRequest represents metadata for user offest pagination.
type usersByOffestPaginationRequest struct {
	Limit int `json:"limit" binding:"min=1"`
	Page  int `json:"page" binding:"min=1"`
}

func (h *UserHandler) GetUsersByOffsetPagination(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		handleError(c, domain.ErrInternalServerError)
		return
	}

	domainToken, ok := token.(*domain.Token)
	if !ok {
		handleError(c, domain.ErrInternalServerError)
	}

	var req usersByOffestPaginationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	users, err := h.userService.GetUsersByOffestPagination(c, domainToken, req.Page, req.Limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapUsersToUserInfoResponse(users))
}

// usersByTimePaginationRequest represents metadata for user time pagination.
type usersByTimePaginationRequest struct {
	After time.Time `json:"after"`
	Limit int       `json:"limit" binding:"min=1"`
}

func (h *UserHandler) GetUsersByTimePagination(c *gin.Context) {
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

	var req usersByTimePaginationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	users, err := h.userService.GetUsersByTimePagination(c, domainToken, req.After, req.Limit)
	if err != nil {
		handleError(c, err)
		return
	}

	userInfo := mapUsersToUserInfoResponse(users)
	if len(userInfo) != 0 {
		cursor := userInfo[len(userInfo)-1].UpdatedAt
		c.JSON(http.StatusOK, gin.H{
			"users":  userInfo,
			"cursor": cursor,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  userInfo,
		"cursor": nil,
	})
}

// searchUserRequest represents a request body for searching a user by username.
type searchUserRequest struct {
	Username string `json:"username" binding:"required"`
	Limit    int    `json:"limit" binding:"min=1"`
}

func (h *UserHandler) SearchUserByUsername(c *gin.Context) {
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

	var req searchUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
	}

	users, err := h.userService.SearchUserByUsername(c, domainToken, req.Username, req.Limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapUsersToUserInfoResponse(users))
}

// searchUserByEmailRequest represents a request body for searching a user by email.
type searchUserByEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Limit int    `json:"limit" binding:"min=1"`
}

func (h *UserHandler) SearchUserByEmail(c *gin.Context) {
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

	var req searchUserByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	users, err := h.userService.SearchUserByEmail(c, domainToken, req.Email, req.Limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapUsersToUserInfoResponse(users))
}

// getUserByIdRequest represents a request body for getting a user by id.
type getUserByIdRequest struct {
	Id uuid.UUID `json:"id" binding:"required"`
}

func (h *UserHandler) GetUserById(c *gin.Context) {
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

	var req getUserByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
	}

	user, err := h.userService.GetUserById(c, domainToken, req.Id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
