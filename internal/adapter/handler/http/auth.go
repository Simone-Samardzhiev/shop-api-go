package http

import (
	"net/http"
	"shop-api-go/internal/adapter/handler/http/request"
	"shop-api-go/internal/adapter/handler/http/response"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService port.AuthService
}

// NewAuthHandler initializes and returns a new AuthHandler instance.
func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary      Authenticate user
// @Description  Authenticates a user using their username and password, returning a new pair of access and refresh tokens upon successful login.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.LoginRequest  true  "Login credentials (username and password)"
// @Success      200      {object}  response.TokensResponse   "Login successful — returns new access and refresh tokens"
// @Failure      400      {object}  response.ErrorResponse    "Invalid request payload or missing fields"
// @Failure      401      {object}  response.ErrorResponse    "Invalid credentials"
// @Failure      500      {object}  response.ErrorResponse    "Internal server error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	tokenGroup, err := h.authService.Login(c, &domain.User{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewTokensResponse(tokenGroup))
}

// RefreshSession godoc
// @Summary      Refresh access token
// @Description  Refreshes the user's authentication session using a valid **refresh token** provided in the Authorization header. Returns a new access/refresh token pair.
// @Tags         Auth
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer refresh token (format: Bearer <token>)"
// @Produce      json
// @Success      200  {object}  response.TokensResponse   "Session refreshed successfully — returns new access and refresh tokens"
// @Failure      401  {object}  response.ErrorResponse    "Missing, malformed, or invalid token"
// @Failure      403  {object}  response.ErrorResponse    "Invalid token type (expected refresh token)"
// @Failure      500  {object}  response.ErrorResponse    "Internal server error"
// @Router       /auth/refresh-session [post]
func (h *AuthHandler) RefreshSession(c *gin.Context) {
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

	tokenGroup, err := h.authService.RefreshSession(c, domainToken)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewTokensResponse(tokenGroup))
}
