package auth_test

import (
	"context"
	"testing"
	"time"

	"hosim-go/internal/auth"
	"hosim-go/internal/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockAuthRepo struct {
	createUserFn             func(ctx context.Context, user *auth.User) (*auth.User, error)
	findByUsernameFn         func(ctx context.Context, username string) (*auth.User, error)
	findByEmailFn            func(ctx context.Context, email string) (*auth.User, error)
	findByUsernameOrEmailFn  func(ctx context.Context, identifier string) (*auth.User, error)
	findUserByIDFn           func(ctx context.Context, id string) (*auth.User, error)
	countUsersFn             func(ctx context.Context) (int64, error)
	saveRefreshTokenFn       func(ctx context.Context, token *auth.RefreshToken) error
	findRefreshTokenFn       func(ctx context.Context, token string) (*auth.RefreshToken, error)
	revokeRefreshTokenFn     func(ctx context.Context, token string) error
	revokeAllUserTokensFn    func(ctx context.Context, userID string) error
	findRoleByCodeFn         func(ctx context.Context, code string) (*auth.Role, error)
	findAllRolesFn           func(ctx context.Context) ([]auth.Role, error)
	findAllPermissionsFn     func(ctx context.Context) ([]auth.Permission, error)
}

func (m *mockAuthRepo) CreateUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, user)
	}
	return user, nil
}

func (m *mockAuthRepo) FindByUsername(ctx context.Context, username string) (*auth.User, error) {
	if m.findByUsernameFn != nil {
		return m.findByUsernameFn(ctx, username)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) FindByUsernameOrEmail(ctx context.Context, identifier string) (*auth.User, error) {
	if m.findByUsernameOrEmailFn != nil {
		return m.findByUsernameOrEmailFn(ctx, identifier)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) FindUserByID(ctx context.Context, id string) (*auth.User, error) {
	if m.findUserByIDFn != nil {
		return m.findUserByIDFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) CountUsers(ctx context.Context) (int64, error) {
	if m.countUsersFn != nil {
		return m.countUsersFn(ctx)
	}
	return 0, nil
}

func (m *mockAuthRepo) SaveRefreshToken(ctx context.Context, token *auth.RefreshToken) error {
	if m.saveRefreshTokenFn != nil {
		return m.saveRefreshTokenFn(ctx, token)
	}
	return nil
}

func (m *mockAuthRepo) FindRefreshToken(ctx context.Context, token string) (*auth.RefreshToken, error) {
	if m.findRefreshTokenFn != nil {
		return m.findRefreshTokenFn(ctx, token)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	if m.revokeRefreshTokenFn != nil {
		return m.revokeRefreshTokenFn(ctx, token)
	}
	return nil
}

func (m *mockAuthRepo) RevokeAllUserTokens(ctx context.Context, userID string) error {
	if m.revokeAllUserTokensFn != nil {
		return m.revokeAllUserTokensFn(ctx, userID)
	}
	return nil
}

func (m *mockAuthRepo) FindRoleByCode(ctx context.Context, code string) (*auth.Role, error) {
	if m.findRoleByCodeFn != nil {
		return m.findRoleByCodeFn(ctx, code)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) FindAllRoles(ctx context.Context) ([]auth.Role, error) {
	if m.findAllRolesFn != nil {
		return m.findAllRolesFn(ctx)
	}
	return []auth.Role{}, nil
}

func (m *mockAuthRepo) FindAllPermissions(ctx context.Context) ([]auth.Permission, error) {
	if m.findAllPermissionsFn != nil {
		return m.findAllPermissionsFn(ctx)
	}
	return []auth.Permission{}, nil
}

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret:          "test-secret-1234567890",
		JWTAccessDuration:  15 * time.Minute,
		JWTRefreshDuration: 7 * 24 * time.Hour,
	}
}

func TestRegister_Success(t *testing.T) {
	doctorRole := &auth.Role{
		ID:   "role-doc",
		Code: "DOCTOR",
		Name: "Dokter",
	}

	repo := &mockAuthRepo{
		findByUsernameFn: func(ctx context.Context, username string) (*auth.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findByEmailFn: func(ctx context.Context, email string) (*auth.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findRoleByCodeFn: func(ctx context.Context, code string) (*auth.Role, error) {
			if code == "DOCTOR" {
				return doctorRole, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
		createUserFn: func(ctx context.Context, user *auth.User) (*auth.User, error) {
			user.ID = "usr-generated-1"
			return user, nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	res, err := svc.Register(context.Background(), auth.RegisterRequest{
		Username: "dokter_andi",
		Email:    "andi@hosim.local",
		Password: "password123",
		Name:     "dr. Andi",
		Role:     "DOCTOR",
	})

	if err != nil {
		t.Fatalf("Register gagal: %v", err)
	}
	if res.Username != "dokter_andi" {
		t.Errorf("Username tidak cocok, didapat: %s", res.Username)
	}
	if res.Role != "DOCTOR" {
		t.Errorf("Role tidak cocok, didapat: %s", res.Role)
	}
}

func TestRegister_UsernameTaken(t *testing.T) {
	repo := &mockAuthRepo{
		findByUsernameFn: func(ctx context.Context, username string) (*auth.User, error) {
			return &auth.User{ID: "existing-1", Username: username}, nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	_, err := svc.Register(context.Background(), auth.RegisterRequest{
		Username: "existing_user",
		Email:    "new@hosim.local",
		Password: "password123",
		Name:     "User",
	})

	if err != auth.ErrUsernameTaken {
		t.Errorf("Diharapkan ErrUsernameTaken, didapat: %v", err)
	}
}

func TestRegister_InvalidRole(t *testing.T) {
	repo := &mockAuthRepo{
		findByUsernameFn: func(ctx context.Context, username string) (*auth.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findByEmailFn: func(ctx context.Context, email string) (*auth.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findRoleByCodeFn: func(ctx context.Context, code string) (*auth.Role, error) {
			return nil, gorm.ErrRecordNotFound // Role tidak ditemukan
		},
	}

	svc := auth.NewService(repo, testConfig())

	_, err := svc.Register(context.Background(), auth.RegisterRequest{
		Username: "new_user",
		Email:    "new@hosim.local",
		Password: "password123",
		Name:     "User Baru",
		Role:     "NON_EXISTENT_ROLE",
	})

	if err == nil {
		t.Fatal("Register dengan role tidak valid seharusnya menghasilkan error, tapi sukses")
	}
}

func TestLogin_Success(t *testing.T) {
	rawPassword := "rahasia123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	adminRole := &auth.Role{
		ID:   "role-admin",
		Code: "ADMIN",
		Name: "Administrator",
		Permissions: []auth.Permission{
			{Code: "department:read", Name: "Lihat Departemen"},
			{Code: "patient:read", Name: "Lihat Pasien"},
		},
	}

	user := &auth.User{
		ID:       "usr-1",
		Username: "admin",
		Email:    "admin@hosim.local",
		Password: string(hashedPassword),
		Role:     adminRole,
		IsActive: true,
	}

	repo := &mockAuthRepo{
		findByUsernameOrEmailFn: func(ctx context.Context, identifier string) (*auth.User, error) {
			return user, nil
		},
		saveRefreshTokenFn: func(ctx context.Context, token *auth.RefreshToken) error {
			return nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	res, err := svc.Login(context.Background(), auth.LoginRequest{
		Username: "admin",
		Password: rawPassword,
	})

	if err != nil {
		t.Fatalf("Login gagal: %v", err)
	}
	if res.AccessToken == "" {
		t.Fatal("AccessToken tidak boleh kosong")
	}
	if res.RefreshToken == "" {
		t.Fatal("RefreshToken tidak boleh kosong")
	}
	if res.User.Username != "admin" {
		t.Errorf("Username tidak sesuai: %s", res.User.Username)
	}
	if res.User.Role != "ADMIN" {
		t.Errorf("Role tidak sesuai: %s", res.User.Role)
	}
	if len(res.User.Permissions) != 2 {
		t.Errorf("Diharapkan 2 permissions, didapat: %d", len(res.User.Permissions))
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	user := &auth.User{
		ID:       "usr-1",
		Username: "admin",
		Password: string(hashedPassword),
		IsActive: true,
	}

	repo := &mockAuthRepo{
		findByUsernameOrEmailFn: func(ctx context.Context, identifier string) (*auth.User, error) {
			return user, nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	_, err := svc.Login(context.Background(), auth.LoginRequest{
		Username: "admin",
		Password: "wrong-password",
	})

	if err != auth.ErrInvalidCredentials {
		t.Errorf("Diharapkan ErrInvalidCredentials, didapat: %v", err)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &auth.User{
		ID:       "usr-inactive",
		Username: "inactive_user",
		Password: string(hashedPassword),
		IsActive: false,
	}

	repo := &mockAuthRepo{
		findByUsernameOrEmailFn: func(ctx context.Context, identifier string) (*auth.User, error) {
			return user, nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	_, err := svc.Login(context.Background(), auth.LoginRequest{
		Username: "inactive_user",
		Password: "password123",
	})

	if err != auth.ErrUserInactive {
		t.Errorf("Diharapkan ErrUserInactive, didapat: %v", err)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	nurseRole := &auth.Role{
		ID:   "role-nurse",
		Code: "NURSE",
		Name: "Perawat",
		Permissions: []auth.Permission{
			{Code: "outpatient:register"},
		},
	}

	user := &auth.User{
		ID:       "usr-1",
		Username: "perawat",
		Role:     nurseRole,
		IsActive: true,
	}

	rt := &auth.RefreshToken{
		ID:        "rt-1",
		UserID:    user.ID,
		User:      user,
		Token:     "valid-refresh-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		RevokedAt: nil,
	}

	revokedToken := ""
	savedToken := ""

	repo := &mockAuthRepo{
		findRefreshTokenFn: func(ctx context.Context, token string) (*auth.RefreshToken, error) {
			return rt, nil
		},
		revokeRefreshTokenFn: func(ctx context.Context, token string) error {
			revokedToken = token
			return nil
		},
		saveRefreshTokenFn: func(ctx context.Context, token *auth.RefreshToken) error {
			savedToken = token.Token
			return nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	res, err := svc.RefreshToken(context.Background(), auth.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	})

	if err != nil {
		t.Fatalf("RefreshToken gagal: %v", err)
	}
	if revokedToken != "valid-refresh-token" {
		t.Errorf("Token lama harus di-revoke, didapat: %s", revokedToken)
	}
	if savedToken == "" {
		t.Error("Token baru harus disimpan ke database")
	}
	if res.AccessToken == "" {
		t.Error("Access token baru tidak boleh kosong")
	}
	if res.User.Role != "NURSE" {
		t.Errorf("Role harus NURSE, didapat: %s", res.User.Role)
	}
}

func TestGetRolesAndPermissions(t *testing.T) {
	repo := &mockAuthRepo{
		findAllRolesFn: func(ctx context.Context) ([]auth.Role, error) {
			return []auth.Role{
				{Code: "ADMIN", Name: "Administrator"},
				{Code: "DOCTOR", Name: "Dokter"},
			}, nil
		},
		findAllPermissionsFn: func(ctx context.Context) ([]auth.Permission, error) {
			return []auth.Permission{
				{Code: "department:read", Name: "Lihat Departemen"},
				{Code: "patient:read", Name: "Lihat Pasien"},
			}, nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	roles, err := svc.GetRoles(context.Background())
	if err != nil {
		t.Fatalf("GetRoles gagal: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("Diharapkan 2 roles, didapat: %d", len(roles))
	}

	perms, err := svc.GetPermissions(context.Background())
	if err != nil {
		t.Fatalf("GetPermissions gagal: %v", err)
	}
	if len(perms) != 2 {
		t.Errorf("Diharapkan 2 permissions, didapat: %d", len(perms))
	}
}
