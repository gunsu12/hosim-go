package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/jwt"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter(secret string) *gin.Engine {
	r := gin.New()

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	{
		api.GET("/profile", func(c *gin.Context) {
			response.Success(c, http.StatusOK, "OK", gin.H{
				"user_id":     middleware.GetUserID(c),
				"username":    middleware.GetUsername(c),
				"user_role":   middleware.GetUserRole(c),
				"permissions": middleware.GetUserPermissions(c),
			})
		})

		api.GET("/admin-only", middleware.RequireRole("ADMIN"), func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Admin Area", nil)
		})

		// Single Permission
		api.DELETE("/departments/:id", middleware.RequirePermission("department:delete"), func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Departemen dihapus", nil)
		})
		api.GET("/departments", middleware.RequirePermission("department:read"), func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Daftar departemen", nil)
		})
		api.POST("/departments", middleware.RequirePermission("department:create"), func(c *gin.Context) {
			response.Success(c, http.StatusCreated, "Departemen dibuat", nil)
		})
		api.POST("/patients", middleware.RequirePermission("patient:create"), func(c *gin.Context) {
			response.Success(c, http.StatusCreated, "Pasien didaftarkan", nil)
		})
		api.DELETE("/patients/:id", middleware.RequirePermission("patient:delete"), func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Pasien dihapus", nil)
		})

		// Multiple Permissions (Opsi 1: OR logic)
		api.GET("/practitioners", middleware.RequireAnyPermission("practitioner:read", "outpatient:register", "appointment:view"), func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Daftar dokter", nil)
		})
	}

	return r
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	router := setupTestRouter("my-secret")

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Diharapkan 401 Unauthorized, didapat: %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	router := setupTestRouter("my-secret")

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "InvalidHeaderFormat")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Diharapkan 401 Unauthorized, didapat: %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "my-secret"
	router := setupTestRouter(secret)

	token, _ := jwt.GenerateAccessToken("usr-99", "budi", "STAFF", []string{"department:read"}, secret, 15*time.Minute)

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Diharapkan 200 OK, didapat: %d (Body: %s)", w.Code, w.Body.String())
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	secret := "my-secret"
	router := setupTestRouter(secret)

	// User adalah DOCTOR, tapi endpoint meminta ADMIN
	token, _ := jwt.GenerateAccessToken("usr-12", "dokter_rudi", "DOCTOR", []string{}, secret, 15*time.Minute)

	req, _ := http.NewRequest(http.MethodGet, "/api/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Diharapkan 403 Forbidden, didapat: %d", w.Code)
	}
}

func TestRequirePermission_Forbidden(t *testing.T) {
	secret := "my-secret"
	router := setupTestRouter(secret)

	// User staf hanya punya permission view, bukan delete
	token, _ := jwt.GenerateAccessToken("usr-2", "staf_biasa", "STAFF", []string{"department:read"}, secret, 15*time.Minute)

	req, _ := http.NewRequest(http.MethodDelete, "/api/departments/123", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Diharapkan 403 Forbidden, didapat: %d", w.Code)
	}
}

func TestRequirePermission_Success(t *testing.T) {
	secret := "my-secret"
	router := setupTestRouter(secret)

	token, _ := jwt.GenerateAccessToken("usr-2", "admin_dept", "STAFF", []string{"department:delete"}, secret, 15*time.Minute)

	req, _ := http.NewRequest(http.MethodDelete, "/api/departments/123", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Diharapkan 200 OK, didapat: %d", w.Code)
	}
}

func TestRequireAnyPermission_ParamedicAccess(t *testing.T) {
	secret := "my-secret"
	router := setupTestRouter(secret)

	// Suster Ani TIDAK punya "practitioner:read", tapi PUNYA "outpatient:register"
	token, _ := jwt.GenerateAccessToken("usr-ani", "suster_ani", "PARAMEDIC", []string{"outpatient:register"}, secret, 15*time.Minute)

	req, _ := http.NewRequest(http.MethodGet, "/api/practitioners", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Diharapkan 200 OK untuk Paramedic dengan izin outpatient:register, didapat: %d", w.Code)
	}
}

func TestRequireAnyPermission_AdminBypass(t *testing.T) {
	secret := "my-secret"
	router := setupTestRouter(secret)

	// Role ADMIN otomatis bypass seluruh permission check
	token, _ := jwt.GenerateAccessToken("usr-admin", "admin", "ADMIN", []string{}, secret, 15*time.Minute)

	req, _ := http.NewRequest(http.MethodGet, "/api/practitioners", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Diharapkan 200 OK untuk Admin bypass, didapat: %d", w.Code)
	}
}

func TestDoctorVsAdminRBAC(t *testing.T) {
	secret := "my-secret"
	router := setupTestRouter(secret)

	// Token DOCTOR dengan permissions standar dari seeder RBAC
	doctorPerms := []string{
		"department:read", "room:read", "service_unit:read", "payer:read", "referal:read", "tariff_class:read",
		"practitioner:read", "patient:read", "patient:create", "patient:update",
		"outpatient:view", "appointment:view",
	}
	doctorToken, _ := jwt.GenerateAccessToken("usr-doc", "dokter", "DOCTOR", doctorPerms, secret, 15*time.Minute)

	// Token ADMIN (Superadmin bypass)
	adminToken, _ := jwt.GenerateAccessToken("usr-adm", "admin", "ADMIN", []string{}, secret, 15*time.Minute)

	// 1. DOCTOR membaca departemen -> Harus Sukses (200)
	req1, _ := http.NewRequest(http.MethodGet, "/api/departments", nil)
	req1.Header.Set("Authorization", "Bearer "+doctorToken)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("Dokter harusnya bisa membaca departemen, didapat: %d", w1.Code)
	}

	// 2. DOCTOR membuat departemen -> Harus Ditolak (403 Forbidden)
	req2, _ := http.NewRequest(http.MethodPost, "/api/departments", nil)
	req2.Header.Set("Authorization", "Bearer "+doctorToken)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusForbidden {
		t.Errorf("Dokter dilarang membuat departemen, diharapkan 403, didapat: %d", w2.Code)
	}

	// 3. DOCTOR mendaftarkan pasien -> Harus Sukses (201 Created)
	req3, _ := http.NewRequest(http.MethodPost, "/api/patients", nil)
	req3.Header.Set("Authorization", "Bearer "+doctorToken)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusCreated {
		t.Errorf("Dokter harusnya bisa mendaftarkan pasien, diharapkan 201, didapat: %d", w3.Code)
	}

	// 4. DOCTOR menghapus pasien -> Harus Ditolak (403 Forbidden)
	req4, _ := http.NewRequest(http.MethodDelete, "/api/patients/123", nil)
	req4.Header.Set("Authorization", "Bearer "+doctorToken)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	if w4.Code != http.StatusForbidden {
		t.Errorf("Dokter dilarang menghapus pasien, diharapkan 403, didapat: %d", w4.Code)
	}

	// 5. ADMIN membuat departemen -> Harus Sukses (201 Created) via Superadmin bypass
	req5, _ := http.NewRequest(http.MethodPost, "/api/departments", nil)
	req5.Header.Set("Authorization", "Bearer "+adminToken)
	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req5)
	if w5.Code != http.StatusCreated {
		t.Errorf("Admin harusnya bisa membuat departemen (bypass), didapat: %d", w5.Code)
	}

	// 6. ADMIN menghapus pasien -> Harus Sukses (200 OK) via Superadmin bypass
	req6, _ := http.NewRequest(http.MethodDelete, "/api/patients/123", nil)
	req6.Header.Set("Authorization", "Bearer "+adminToken)
	w6 := httptest.NewRecorder()
	router.ServeHTTP(w6, req6)
	if w6.Code != http.StatusOK {
		t.Errorf("Admin harusnya bisa menghapus pasien (bypass), didapat: %d", w6.Code)
	}
}
