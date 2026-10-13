package auth_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"hosim-go/internal/auth"
	"hosim-go/internal/config"
	"hosim-go/internal/practitioner"
	"hosim-go/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockAuthRepo struct {
	createUserFn            func(ctx context.Context, user *auth.User) (*auth.User, error)
	findByUsernameFn        func(ctx context.Context, username string) (*auth.User, error)
	findByEmailFn           func(ctx context.Context, email string) (*auth.User, error)
	findByUsernameOrEmailFn func(ctx context.Context, identifier string) (*auth.User, error)
	findUserByIDFn          func(ctx context.Context, id string) (*auth.User, error)
	findByPractitionerIDFn  func(ctx context.Context, practitionerID string) (*auth.User, error)
	countUsersFn            func(ctx context.Context) (int64, error)
	updateUserFn            func(ctx context.Context, user *auth.User) error
	saveRefreshTokenFn      func(ctx context.Context, token *auth.RefreshToken) error
	findRefreshTokenFn      func(ctx context.Context, token string) (*auth.RefreshToken, error)
	revokeRefreshTokenFn    func(ctx context.Context, token string) error
	revokeAllUserTokensFn   func(ctx context.Context, userID string) error
	findRoleByCodeFn        func(ctx context.Context, code string) (*auth.Role, error)
	findAllRolesFn          func(ctx context.Context) ([]auth.Role, error)
	findAllPermissionsFn    func(ctx context.Context) ([]auth.Permission, error)
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

func (m *mockAuthRepo) FindByPractitionerID(ctx context.Context, practitionerID string) (*auth.User, error) {
	if m.findByPractitionerIDFn != nil {
		return m.findByPractitionerIDFn(ctx, practitionerID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) UpdateUser(ctx context.Context, user *auth.User) error {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, user)
	}
	return nil
}

type mockPractitionerReader struct {
	getPractitionerByIDFn func(ctx context.Context, id string) (*practitioner.Practitioner, error)
}

func (m *mockPractitionerReader) GetPractitionerByID(ctx context.Context, id string) (*practitioner.Practitioner, error) {
	if m.getPractitionerByIDFn != nil {
		return m.getPractitionerByIDFn(ctx, id)
	}
	return nil, errors.New("practitioner not found")
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

func (m *mockAuthRepo) ListUsers(ctx context.Context, search string, roleID string, isActive *bool, page, limit int) ([]auth.User, int64, error) {
	return []auth.User{}, 0, nil
}

func (m *mockAuthRepo) DeleteUser(ctx context.Context, id string, deletedBy string) error {
	return nil
}

func (m *mockAuthRepo) FindRoleByID(ctx context.Context, id string) (*auth.Role, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) CreateRole(ctx context.Context, role *auth.Role, permIDs []string) (*auth.Role, error) {
	return role, nil
}

func (m *mockAuthRepo) UpdateRole(ctx context.Context, role *auth.Role, permIDs []string) (*auth.Role, error) {
	return role, nil
}

func (m *mockAuthRepo) DeleteRole(ctx context.Context, id string) error {
	return nil
}

func (m *mockAuthRepo) FindAllPermissions(ctx context.Context) ([]auth.Permission, error) {
	if m.findAllPermissionsFn != nil {
		return m.findAllPermissionsFn(ctx)
	}
	return []auth.Permission{}, nil
}

func (m *mockAuthRepo) FindPermissionByID(ctx context.Context, id string) (*auth.Permission, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) FindPermissionByCode(ctx context.Context, code string) (*auth.Permission, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) CreatePermission(ctx context.Context, perm *auth.Permission) (*auth.Permission, error) {
	return perm, nil
}

func (m *mockAuthRepo) UpdatePermission(ctx context.Context, perm *auth.Permission) error {
	return nil
}

func (m *mockAuthRepo) DeletePermission(ctx context.Context, id string) error {
	return nil
}

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret:          "test-secret-1234567890",
		JWTAccessDuration:  15 * time.Minute,
		JWTRefreshDuration: 7 * 24 * time.Hour,
	}
}

func TestRegister_Success(t *testing.T) {
	staffRole := &auth.Role{
		ID:   "role-staff",
		Code: "STAFF",
		Name: "Staf",
	}

	repo := &mockAuthRepo{
		findByUsernameFn: func(ctx context.Context, username string) (*auth.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findByEmailFn: func(ctx context.Context, email string) (*auth.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findRoleByCodeFn: func(ctx context.Context, code string) (*auth.Role, error) {
			if code == "STAFF" {
				return staffRole, nil
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
		Username: "staf_andi",
		Email:    "andi@hosim.local",
		Password: "password123",
		Name:     "Andi",
		Role:     "STAFF",
	})

	if err != nil {
		t.Fatalf("Register gagal: %v", err)
	}
	if res.Username != "staf_andi" {
		t.Errorf("Username tidak cocok, didapat: %s", res.Username)
	}
	if res.Role != "STAFF" {
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

func TestLinkAndUnlinkPractitioner(t *testing.T) {
	user := &auth.User{
		ID:       "usr-1",
		Username: "dokter_andi",
		Name:     "dr. Andi",
		IsActive: true,
	}

	activePrac := &practitioner.Practitioner{
		ID:       "prac-1",
		Name:     "dr. Andi Sp.A",
		IsActive: true,
	}

	revokedUserTokens := ""
	var updatedUser *auth.User

	repo := &mockAuthRepo{
		findUserByIDFn: func(ctx context.Context, id string) (*auth.User, error) {
			if id == "usr-1" {
				return user, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
		findByPractitionerIDFn: func(ctx context.Context, practitionerID string) (*auth.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		updateUserFn: func(ctx context.Context, u *auth.User) error {
			updatedUser = u
			return nil
		},
		revokeAllUserTokensFn: func(ctx context.Context, userID string) error {
			revokedUserTokens = userID
			return nil
		},
	}

	pracReader := &mockPractitionerReader{
		getPractitionerByIDFn: func(ctx context.Context, id string) (*practitioner.Practitioner, error) {
			if id == "prac-1" {
				return activePrac, nil
			}
			return nil, auth.ErrPractitionerNotFound
		},
	}

	svc := auth.NewService(repo, testConfig(), pracReader)

	// 1. Link Practitioner
	res, err := svc.LinkPractitioner(context.Background(), "usr-1", "prac-1", "operator-admin")
	if err != nil {
		t.Fatalf("LinkPractitioner gagal: %v", err)
	}
	if res.PractitionerID == nil || *res.PractitionerID != "prac-1" {
		t.Errorf("PractitionerID harus terisi 'prac-1', didapat: %v", res.PractitionerID)
	}
	if updatedUser == nil || updatedUser.PractitionerID == nil || *updatedUser.PractitionerID != "prac-1" {
		t.Errorf("Updated user model harus memiliki PractitionerID 'prac-1'")
	}
	if revokedUserTokens != "usr-1" {
		t.Errorf("RevokeAllUserTokens harus dipanggil untuk usr-1, didapat: %s", revokedUserTokens)
	}

	// 2. Unlink Practitioner
	resUnlink, err := svc.UnlinkPractitioner(context.Background(), "usr-1", "operator-admin")
	if err != nil {
		t.Fatalf("UnlinkPractitioner gagal: %v", err)
	}
	if resUnlink.PractitionerID != nil {
		t.Errorf("PractitionerID harus nil setelah unlink, didapat: %v", *resUnlink.PractitionerID)
	}
	if updatedUser.PractitionerID != nil {
		t.Errorf("Updated user model PractitionerID harus nil setelah unlink")
	}
}

func TestLinkPractitioner_Validation(t *testing.T) {
	user := &auth.User{
		ID:       "usr-1",
		Username: "dokter_andi",
		IsActive: true,
	}

	otherUser := &auth.User{
		ID:       "usr-other",
		Username: "dokter_lain",
	}

	inactivePrac := &practitioner.Practitioner{
		ID:       "prac-inactive",
		Name:     "dr. Nonaktif",
		IsActive: false,
	}

	activePrac := &practitioner.Practitioner{
		ID:       "prac-active",
		Name:     "dr. Aktif",
		IsActive: true,
	}

	repo := &mockAuthRepo{
		findUserByIDFn: func(ctx context.Context, id string) (*auth.User, error) {
			if id == "usr-1" {
				return user, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
		findByPractitionerIDFn: func(ctx context.Context, practitionerID string) (*auth.User, error) {
			if practitionerID == "prac-active" {
				return otherUser, nil // sudah di-link ke otherUser
			}
			return nil, gorm.ErrRecordNotFound
		},
	}

	pracReader := &mockPractitionerReader{
		getPractitionerByIDFn: func(ctx context.Context, id string) (*practitioner.Practitioner, error) {
			if id == "prac-inactive" {
				return inactivePrac, nil
			}
			if id == "prac-active" {
				return activePrac, nil
			}
			return nil, auth.ErrPractitionerNotFound
		},
	}

	svc := auth.NewService(repo, testConfig(), pracReader)

	// Kasus 1: Nakes tidak ditemukan
	_, err := svc.LinkPractitioner(context.Background(), "usr-1", "prac-unknown", "admin")
	if !errors.Is(err, auth.ErrPractitionerNotFound) {
		t.Errorf("Diharapkan ErrPractitionerNotFound, didapat: %v", err)
	}

	// Kasus 2: Nakes non-aktif
	_, err = svc.LinkPractitioner(context.Background(), "usr-1", "prac-inactive", "admin")
	if !errors.Is(err, auth.ErrPractitionerInactive) {
		t.Errorf("Diharapkan ErrPractitionerInactive, didapat: %v", err)
	}

	// Kasus 3: Nakes sudah terhubung ke user lain
	_, err = svc.LinkPractitioner(context.Background(), "usr-1", "prac-active", "admin")
	if !errors.Is(err, auth.ErrPractitionerAlreadyLinked) {
		t.Errorf("Diharapkan ErrPractitionerAlreadyLinked, didapat: %v", err)
	}
}

func TestLogin_WithPractitionerID(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("rahasia123"), bcrypt.DefaultCost)
	pracID := "prac-777"
	role := &auth.Role{ID: "r-doc", Code: "DOCTOR", Name: "Dokter"}

	user := &auth.User{
		ID:             "usr-doc-1",
		Username:       "dr_budi",
		Email:          "budi@hospital.com",
		Password:       string(hashedPassword),
		RoleID:         &role.ID,
		Role:           role,
		PractitionerID: &pracID,
		IsActive:       true,
	}

	repo := &mockAuthRepo{
		findByUsernameOrEmailFn: func(ctx context.Context, identifier string) (*auth.User, error) {
			return user, nil
		},
		saveRefreshTokenFn: func(ctx context.Context, token *auth.RefreshToken) error {
			return nil
		},
	}

	cfg := testConfig()
	svc := auth.NewService(repo, cfg)

	res, err := svc.Login(context.Background(), auth.LoginRequest{
		Username: "dr_budi",
		Password: "rahasia123",
	})

	if err != nil {
		t.Fatalf("Login gagal: %v", err)
	}

	if res.User.PractitionerID == nil || *res.User.PractitionerID != pracID {
		t.Errorf("UserResponse PractitionerID harus %s, didapat: %v", pracID, res.User.PractitionerID)
	}

	// Verifikasi klaim di Access Token
	claims, err := jwt.ValidateToken(res.AccessToken, cfg.JWTSecret)
	if err != nil {
		t.Fatalf("ValidateToken gagal: %v", err)
	}

	if claims.PractitionerID == nil || *claims.PractitionerID != pracID {
		t.Errorf("Klaim JWT PractitionerID harus %s, didapat: %v", pracID, claims.PractitionerID)
	}
}

func TestRegister_CannotSelfRegisterAdmin(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := auth.NewService(repo, testConfig())

	_, err := svc.Register(context.Background(), auth.RegisterRequest{
		Username: "new_admin",
		Email:    "new_admin@hospital.com",
		Password: "password123",
		Name:     "Hacker",
		Role:     "ADMIN",
	})

	if !errors.Is(err, auth.ErrSelfAdminRegister) {
		t.Errorf("Pendaftaran mandiri dengan role ADMIN harus ditolak dengan ErrSelfAdminRegister, didapat: %v", err)
	}
}

func TestRegister_CannotSelfRegisterPrivilegedRoles(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := auth.NewService(repo, testConfig())

	for _, role := range []string{"DOCTOR", "PARAMEDIC", "NURSE", "PHARMACIST"} {
		_, err := svc.Register(context.Background(), auth.RegisterRequest{
			Username: "new_" + strings.ToLower(role),
			Email:    "new_" + strings.ToLower(role) + "@hospital.com",
			Password: "password123",
			Name:     "Test " + role,
			Role:     role,
		})

		if !errors.Is(err, auth.ErrSelfRegisterPrivilegedRole) {
			t.Errorf("Pendaftaran mandiri role %s harus ditolak dengan ErrSelfRegisterPrivilegedRole, didapat: %v", role, err)
		}
	}
}

func TestUpdateUser_RevokesTokensOnSecurityChanges(t *testing.T) {
	user := &auth.User{
		ID:       "usr-target",
		Username: "staff_1",
		Email:    "staff@hospital.com",
		IsActive: true,
	}

	var revokedUser string
	repo := &mockAuthRepo{
		findUserByIDFn: func(ctx context.Context, id string) (*auth.User, error) {
			if id == "usr-target" {
				return user, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
		updateUserFn: func(ctx context.Context, u *auth.User) error {
			return nil
		},
		revokeAllUserTokensFn: func(ctx context.Context, userID string) error {
			revokedUser = userID
			return nil
		},
	}

	svc := auth.NewService(repo, testConfig())

	// 1. Password reset memicu token revocation
	revokedUser = ""
	newPass := "NewPass123!"
	_, err := svc.UpdateUser(context.Background(), "usr-target", auth.UpdateUserRequest{
		Name:     "Staff Updated",
		Email:    "staff@hospital.com",
		Password: newPass,
	}, "admin")
	if err != nil {
		t.Fatalf("UpdateUser gagal: %v", err)
	}
	if revokedUser != "usr-target" {
		t.Errorf("Password reset harus memicu RevokeAllUserTokens")
	}

	// 2. Deaktivasi memicu token revocation
	revokedUser = ""
	inactive := false
	_, err = svc.UpdateUser(context.Background(), "usr-target", auth.UpdateUserRequest{
		Name:     "Staff Updated",
		Email:    "staff@hospital.com",
		IsActive: &inactive,
	}, "admin")
	if err != nil {
		t.Fatalf("UpdateUser gagal: %v", err)
	}
	if revokedUser != "usr-target" {
		t.Errorf("Deaktivasi akun harus memicu RevokeAllUserTokens")
	}
}
