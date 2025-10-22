package request

// LoginRequest represent a login request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
