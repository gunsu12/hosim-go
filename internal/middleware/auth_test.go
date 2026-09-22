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
