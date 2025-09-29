package http

import (
	"net/http"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"

	"github.com/gin-gonic/gin"
)

// AuthHandler represent HTTP handler for auth-related requests.
type AuthHandler struct {
	authService port.AuthService
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// loginRequest represents a request body for logging in.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindingError(c, err)
		return
	}

	tokenGroup, err := h.authService.Login(c, &domain.User{Username: req.Username, Password: req.Password})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  tokenGroup.AccessToken,
		"refreshToken": tokenGroup.RefreshToken,
	})
}
