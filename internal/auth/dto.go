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
	RoleID      *string   `json:"role_id,omitempty"`
	Role        string    `json:"role"`
	RoleName    string    `json:"role_name,omitempty"`
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
	roleCode := user.GetRoleCode()
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}
	return UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Name:        user.Name,
		RoleID:      user.RoleID,
		Role:        roleCode,
		RoleName:    roleName,
		Permissions: user.GetPermissionCodes(),
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
	}
}

// ========================================================
// DTO Manajemen Pengguna (Users)
// ========================================================

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	RoleID   *string `json:"role_id"`
	IsActive *bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	Email    string  `json:"email" binding:"required,email"`
	RoleID   *string `json:"role_id"`
	IsActive *bool   `json:"is_active"`
	Password string  `json:"password,omitempty" binding:"omitempty,min=6"` // Opsional, diisi jika ingin reset password
}

// ========================================================
// DTO Manajemen Peran (Roles)
// ========================================================

type CreateRoleRequest struct {
	Code          string   `json:"code" binding:"required,min=2,max=50"`
	Name          string   `json:"name" binding:"required,min=2,max=150"`
	Description   string   `json:"description,omitempty"`
	PermissionIDs []string `json:"permission_ids,omitempty"`
}

type UpdateRoleRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=150"`
	Description   string   `json:"description,omitempty"`
	PermissionIDs []string `json:"permission_ids,omitempty"`
}

// ========================================================
// DTO Manajemen Hak Akses (Permissions)
// ========================================================

type CreatePermissionRequest struct {
	Code        string `json:"code" binding:"required,min=3,max=100"`
	Name        string `json:"name" binding:"required,min=2,max=150"`
	Module      string `json:"module" binding:"required,min=2,max=100"`
	Description string `json:"description,omitempty"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=150"`
	Module      string `json:"module" binding:"required,min=2,max=100"`
	Description string `json:"description,omitempty"`
}
