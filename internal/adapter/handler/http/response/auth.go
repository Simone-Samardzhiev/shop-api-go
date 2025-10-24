package response

import "shop-api-go/internal/core/domain"

// TokensResponse represents tokens response.
type TokensResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// NewTokensResponse creates a new TokensResponse instance.
func NewTokensResponse(group *domain.TokenGroup) *TokensResponse {
	return &TokensResponse{
		AccessToken:  group.AccessToken,
		RefreshToken: group.RefreshToken,
	}
}
