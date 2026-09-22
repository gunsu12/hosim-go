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
	ErrInvalidCredentials  = errors.New("username atau password tidak valid")
	ErrUserInactive        = errors.New("akun pengguna telah dinonaktifkan")
	ErrUsernameTaken       = errors.New("username sudah digunakan")
	ErrEmailTaken          = errors.New("email sudah digunakan")
	ErrInvalidRefreshToken = errors.New("refresh token tidak valid atau telah kedaluwarsa")
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context, userID string) (*UserResponse, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetPermissions(ctx context.Context) ([]Permission, error)
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

	// Cek apakah username sudah dipakai
	if existingUser, err := s.repo.FindByUsername(ctx, req.Username); err == nil && existingUser != nil {
		return nil, ErrUsernameTaken
	}

	// Cek apakah email sudah dipakai
	if existingUser, err := s.repo.FindByEmail(ctx, req.Email); err == nil && existingUser != nil {
		return nil, ErrEmailTaken
	}

	// Hash password menggunakan bcrypt
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

	// Verifikasi kata sandi
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	roleCode := user.GetRoleCode()
	if roleCode == "" {
		roleCode = "STAFF"
	}
	permissions := user.GetPermissionCodes()

	// Buat JWT Access Token yang membawa claims role & permissions
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Username, roleCode, permissions, s.cfg.JWTSecret, s.cfg.JWTAccessDuration)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token: %w", err)
	}

	// Buat Refresh Token acak yang aman
	refreshTokenStr, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat refresh token: %w", err)
	}

	// Simpan Refresh Token ke Database
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

	// Revoke refresh token lama (Token Rotation)
	_ = s.repo.RevokeRefreshToken(ctx, rt.Token)

	roleCode := user.GetRoleCode()
	if roleCode == "" {
		roleCode = "STAFF"
	}
	permissions := user.GetPermissionCodes()

	// Buat Access Token baru
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Username, roleCode, permissions, s.cfg.JWTSecret, s.cfg.JWTAccessDuration)
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui access token: %w", err)
	}

	// Buat Refresh Token baru
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

func (s *service) GetRoles(ctx context.Context) ([]Role, error) {
	return s.repo.FindAllRoles(ctx)
}

func (s *service) GetPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.FindAllPermissions(ctx)
}
