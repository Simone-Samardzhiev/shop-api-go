package response

import "shop-api-go/internal/core/domain"

// TokensResponse represents tokens response.
type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewTokensResponse creates a new TokensResponse instance.
func NewTokensResponse(group *domain.TokenGroup) *TokensResponse {
	return &TokensResponse{
		AccessToken:  group.AccessToken,
		RefreshToken: group.RefreshToken,
	}
}
