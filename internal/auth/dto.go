package auth

import (
	"time"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Role     string `json:"role" binding:"omitempty,oneof=ADMIN DOCTOR NURSE PHARMACIST CASHIER STAFF"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"` // Dapat diisi Username atau Email
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserResponse struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"` // dalam satuan detik
	User         UserResponse `json:"user"`
}

func ToUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.GetRoleCode(),
		Permissions: user.GetPermissionCodes(),
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
	}
}
