package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"hosim-go/internal/config"
	"hosim-go/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound        = errors.New("pengguna tidak ditemukan")
	ErrRoleNotFound        = errors.New("peran (role) tidak ditemukan")
	ErrPermissionNotFound  = errors.New("hak akses (permission) tidak ditemukan")
	ErrInvalidCredentials  = errors.New("username atau password tidak valid")
	ErrUserInactive        = errors.New("akun pengguna telah dinonaktifkan")
	ErrUsernameTaken       = errors.New("username sudah digunakan")
	ErrEmailTaken          = errors.New("email sudah digunakan")
	ErrRoleCodeTaken       = errors.New("kode peran (role code) sudah digunakan")
	ErrPermCodeTaken       = errors.New("kode hak akses (permission code) sudah digunakan")
	ErrCannotDeleteAdmin   = errors.New("peran ADMIN adalah peran sistem bawaan dan tidak dapat dihapus")
	ErrInvalidRefreshToken = errors.New("refresh token tidak valid atau telah kedaluwarsa")
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context, userID string) (*UserResponse, error)

	// User Management
	ListUsers(ctx context.Context, search string, roleID string, isActive *bool, page, limit int) ([]UserResponse, int64, error)
	CreateUser(ctx context.Context, req CreateUserRequest, operator string) (*UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*UserResponse, error)
	UpdateUser(ctx context.Context, id string, req UpdateUserRequest, operator string) (*UserResponse, error)
	DeleteUser(ctx context.Context, id string, operator string) error

	// Role Management
	GetRoles(ctx context.Context) ([]Role, error)
	GetRoleByID(ctx context.Context, id string) (*Role, error)
	CreateRole(ctx context.Context, req CreateRoleRequest) (*Role, error)
	UpdateRole(ctx context.Context, id string, req UpdateRoleRequest) (*Role, error)
	DeleteRole(ctx context.Context, id string) error

	// Permission Management
	GetPermissions(ctx context.Context) ([]Permission, error)
	CreatePermission(ctx context.Context, req CreatePermissionRequest) (*Permission, error)
	UpdatePermission(ctx context.Context, id string, req UpdatePermissionRequest) (*Permission, error)
	DeletePermission(ctx context.Context, id string) error
}

type service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*UserResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if existingUser, err := s.repo.FindByUsername(ctx, req.Username); err == nil && existingUser != nil {
		return nil, ErrUsernameTaken
	}

	if existingUser, err := s.repo.FindByEmail(ctx, req.Email); err == nil && existingUser != nil {
		return nil, ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi kata sandi: %w", err)
	}

	roleName := strings.ToUpper(strings.TrimSpace(req.Role))
	if roleName == "" {
		roleName = "STAFF"
	}

	roleModel, err := s.repo.FindRoleByCode(ctx, roleName)
	if err != nil || roleModel == nil {
		return nil, fmt.Errorf("role '%s' tidak valid atau belum terdaftar dalam sistem", roleName)
	}

	newUser := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     strings.TrimSpace(req.Name),
		RoleID:   &roleModel.ID,
		Role:     roleModel,
		IsActive: true,
	}

	created, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat akun pengguna: %w", err)
	}

	res := ToUserResponse(created)
	return &res, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	identifier := strings.TrimSpace(req.Username)
	user, err := s.repo.FindByUsernameOrEmail(ctx, identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	roleCode := user.GetRoleCode()
	if roleCode == "" {
		roleCode = "STAFF"
	}
	permissions := user.GetPermissionCodes()

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Username, roleCode, permissions, s.cfg.JWTSecret, s.cfg.JWTAccessDuration)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token: %w", err)
	}

	refreshTokenStr, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat refresh token: %w", err)
	}

	refreshTokenModel := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(s.cfg.JWTRefreshDuration),
	}

	if err := s.repo.SaveRefreshToken(ctx, refreshTokenModel); err != nil {
		return nil, fmt.Errorf("gagal menyimpan refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWTAccessDuration.Seconds()),
		User:         ToUserResponse(user),
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error) {
	rt, err := s.repo.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	if !rt.IsValid() {
		return nil, ErrInvalidRefreshToken
	}

	user := rt.User
	if user == nil {
		user, err = s.repo.FindUserByID(ctx, rt.UserID)
		if err != nil {
			return nil, ErrUserNotFound
		}
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	_ = s.repo.RevokeRefreshToken(ctx, rt.Token)

	roleCode := user.GetRoleCode()
	if roleCode == "" {
		roleCode = "STAFF"
	}
	permissions := user.GetPermissionCodes()

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Username, roleCode, permissions, s.cfg.JWTSecret, s.cfg.JWTAccessDuration)
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui access token: %w", err)
	}

	newRefreshTokenStr, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat refresh token baru: %w", err)
	}

	newRefreshTokenModel := &RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshTokenStr,
		ExpiresAt: time.Now().Add(s.cfg.JWTRefreshDuration),
	}

	if err := s.repo.SaveRefreshToken(ctx, newRefreshTokenModel); err != nil {
		return nil, fmt.Errorf("gagal menyimpan refresh token baru: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWTAccessDuration.Seconds()),
		User:         ToUserResponse(user),
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *service) GetMe(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	res := ToUserResponse(user)
	return &res, nil
}

// ========================================================
// User Management Implementations
// ========================================================

func (s *service) ListUsers(ctx context.Context, search string, roleID string, isActive *bool, page, limit int) ([]UserResponse, int64, error) {
	users, total, err := s.repo.ListUsers(ctx, search, roleID, isActive, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]UserResponse, len(users))
	for i, u := range users {
		userCopy := u
		responses[i] = ToUserResponse(&userCopy)
	}
	return responses, total, nil
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest, operator string) (*UserResponse, error) {
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(strings.ToLower(req.Email))

	if existing, err := s.repo.FindByUsername(ctx, username); err == nil && existing != nil {
		return nil, ErrUsernameTaken
	}
	if existing, err := s.repo.FindByEmail(ctx, email); err == nil && existing != nil {
		return nil, ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi kata sandi: %w", err)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	newUser := &User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Name:      strings.TrimSpace(req.Name),
		RoleID:    req.RoleID,
		IsActive:  isActive,
		CreatedBy: operator,
		UpdatedBy: operator,
	}

	created, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	fullUser, err := s.repo.FindUserByID(ctx, created.ID)
	if err != nil {
		res := ToUserResponse(created)
		return &res, nil
	}

	res := ToUserResponse(fullUser)
	return &res, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	res := ToUserResponse(user)
	return &res, nil
}

func (s *service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest, operator string) (*UserResponse, error) {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if !strings.EqualFold(user.Email, email) {
		if existing, err := s.repo.FindByEmail(ctx, email); err == nil && existing != nil && existing.ID != id {
			return nil, ErrEmailTaken
		}
	}

	user.Name = strings.TrimSpace(req.Name)
	user.Email = email
	user.RoleID = req.RoleID
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	user.UpdatedBy = operator

	// Jika ada reset password
	if strings.TrimSpace(req.Password) != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("gagal mengenkripsi kata sandi baru: %w", err)
		}
		user.Password = string(hashed)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	fullUser, err := s.repo.FindUserByID(ctx, user.ID)
	if err != nil {
		res := ToUserResponse(user)
		return &res, nil
	}

	res := ToUserResponse(fullUser)
	return &res, nil
}

func (s *service) DeleteUser(ctx context.Context, id string, operator string) error {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if strings.EqualFold(user.Username, "admin") {
		return errors.New("pengguna admin utama tidak dapat dihapus")
	}

	if err := s.repo.DeleteUser(ctx, id, operator); err != nil {
		return err
	}
	_ = s.repo.RevokeAllUserTokens(ctx, id)
	return nil
}

// ========================================================
// Role Management Implementations
// ========================================================

func (s *service) GetRoles(ctx context.Context) ([]Role, error) {
	return s.repo.FindAllRoles(ctx)
}

func (s *service) GetRoleByID(ctx context.Context, id string) (*Role, error) {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return role, nil
}

func (s *service) CreateRole(ctx context.Context, req CreateRoleRequest) (*Role, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if existing, err := s.repo.FindRoleByCode(ctx, code); err == nil && existing != nil {
		return nil, ErrRoleCodeTaken
	}

	newRole := &Role{
		Code:        code,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	}

	return s.repo.CreateRole(ctx, newRole, req.PermissionIDs)
}

func (s *service) UpdateRole(ctx context.Context, id string, req UpdateRoleRequest) (*Role, error) {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	role.Name = strings.TrimSpace(req.Name)
	role.Description = strings.TrimSpace(req.Description)

	return s.repo.UpdateRole(ctx, role, req.PermissionIDs)
}

func (s *service) DeleteRole(ctx context.Context, id string) error {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	if strings.EqualFold(role.Code, "ADMIN") {
		return ErrCannotDeleteAdmin
	}

	return s.repo.DeleteRole(ctx, id)
}

// ========================================================
// Permission Management Implementations
// ========================================================

func (s *service) GetPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.FindAllPermissions(ctx)
}

func (s *service) CreatePermission(ctx context.Context, req CreatePermissionRequest) (*Permission, error) {
	code := strings.TrimSpace(strings.ToLower(req.Code))
	if existing, err := s.repo.FindPermissionByCode(ctx, code); err == nil && existing != nil {
		return nil, ErrPermCodeTaken
	}

	newPerm := &Permission{
		Code:        code,
		Name:        strings.TrimSpace(req.Name),
		Module:      strings.TrimSpace(req.Module),
		Description: strings.TrimSpace(req.Description),
	}

	return s.repo.CreatePermission(ctx, newPerm)
}

func (s *service) UpdatePermission(ctx context.Context, id string, req UpdatePermissionRequest) (*Permission, error) {
	perm, err := s.repo.FindPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	perm.Name = strings.TrimSpace(req.Name)
	perm.Module = strings.TrimSpace(req.Module)
	perm.Description = strings.TrimSpace(req.Description)

	if err := s.repo.UpdatePermission(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (s *service) DeletePermission(ctx context.Context, id string) error {
	_, err := s.repo.FindPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPermissionNotFound
		}
		return err
	}
	return s.repo.DeletePermission(ctx, id)
}
