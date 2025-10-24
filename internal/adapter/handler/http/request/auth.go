package request

// LoginRequest represent a login request body.
type LoginRequest struct {
	Username string `json:"username" example:"MyUsername"`
	Password string `json:"password" example:"Secret_password123"`
}
