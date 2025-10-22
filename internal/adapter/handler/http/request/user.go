package request

// RegisterRequest represents a register request.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,min_bytes=8,max_bytes=255"`
	Username string `json:"username" binding:"required,min_bytes=8,max_bytes=255"`
	Password string `json:"password" binding:"required,password"`
}

// UpdateAccountRequest represents a update to user's account.
type UpdateAccountRequest struct {
	Username    string  `json:"username" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	NewUsername *string `json:"newUsername" binding:"omitempty,min_bytes=8,max_bytes=255"`
	NewEmail    *string `json:"newEmail" binding:"omitempty,min_bytes=8,max_bytes=255"`
	NewPassword *string `json:"newPassword" binding:"omitempty,password"`
}
