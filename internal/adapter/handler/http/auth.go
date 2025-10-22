package http

import (
	"net/http"
	"shop-api-go/internal/adapter/handler/http/request"
	"shop-api-go/internal/adapter/handler/http/response"
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	tokenGroup, err := h.authService.Login(c, &domain.User{Username: req.Username, Password: req.Password})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewTokensResponse(tokenGroup))
}

func (h *AuthHandler) RefreshSession(c *gin.Context) {
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

	tokenGroup, err := h.authService.RefreshSession(c, domainToken)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewTokensResponse(tokenGroup))
}
