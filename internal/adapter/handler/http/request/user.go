package request

// RegisterRequest represents a register request body.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,min_bytes=8,max_bytes=255" example:"newUser@email.com"`
	Username string `json:"username" binding:"required,min_bytes=8,max_bytes=255" example:"newUser123"`
	Password string `json:"password" binding:"required,password" example:"NewSecret_123"`
}

// UpdateAccountRequest represents a for update to user's account request body.
type UpdateAccountRequest struct {
	Username    string  `json:"username" binding:"required" example:"MyUsername"`
	Password    string  `json:"password" binding:"required" example:"MyPassword_123"`
	NewUsername *string `json:"newUsername" binding:"omitempty,min_bytes=8,max_bytes=255" example:"newUsername123"`
	NewEmail    *string `json:"newEmail" binding:"omitempty,min_bytes=8,max_bytes=255" example:"newEmail@emai.com"`
	NewPassword *string `json:"newPassword" binding:"omitempty,password" example:"NewSecret_123"`
}
