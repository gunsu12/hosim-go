package middleware

import (
	"net/http"
	"strings"

	"hosim-go/pkg/jwt"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID          = "user_id"
	ContextUsername        = "username"
	ContextUserRole        = "user_role"
	ContextUserPermissions = "user_permissions"
)

// AuthMiddleware memvalidasi JWT Bearer token pada request header
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Header Authorization diperlukan", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "Format Authorization header harus 'Bearer <token>'", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Autentikasi gagal: "+err.Error(), nil)
			c.Abort()
			return
		}

		// Simpan informasi user ke context Gin
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextUserRole, claims.Role)
		c.Set(ContextUserPermissions, claims.Permissions)

		c.Next()
	}
}

// RequireRole membatasi akses endpoint hanya untuk role tertentu
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(ContextUserRole)
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Identitas pengguna tidak ditemukan pada sesi", nil)
			c.Abort()
			return
		}

		currentRole, ok := userRole.(string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "Gagal membaca role pengguna", nil)
			c.Abort()
			return
		}

		// Role ADMIN selalu memiliki akses penuh (Superadmin bypass)
		if strings.EqualFold(currentRole, "ADMIN") {
			c.Next()
			return
		}

		for _, role := range allowedRoles {
			if strings.EqualFold(currentRole, role) {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki wewenang untuk mengakses sumber daya ini", nil)
		c.Abort()
	}
}

// RequirePermission membatasi akses endpoint hanya untuk 1 izin spesifik
func RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if HasPermission(c, permissionCode) {
			c.Next()
			return
		}

		response.Error(c, http.StatusForbidden, "Akses ditolak: Memerlukan izin '"+permissionCode+"'", nil)
		c.Abort()
	}
}

// RequireAnyPermission (Opsi 1) membatasi akses endpoint agar bisa diakses oleh SALAH SATU dari izin yang ditentukan (OR logic)
func RequireAnyPermission(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if HasAnyPermission(c, permissionCodes...) {
			c.Next()
			return
		}

		response.Error(c, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki salah satu izin yang diperlukan untuk aksi ini", nil)
		c.Abort()
	}
}

// HasPermission mengecek apakah user di context memiliki permission tertentu (atau role ADMIN)
func HasPermission(c *gin.Context, permissionCode string) bool {
	// Superadmin bypass
	if strings.EqualFold(GetUserRole(c), "ADMIN") {
		return true
	}

	userPermissions := GetUserPermissions(c)
	for _, p := range userPermissions {
		if p == "*" || strings.EqualFold(p, permissionCode) {
			return true
		}
	}
	return false
}

// HasAnyPermission mengecek apakah user memiliki setidaknya salah satu permission yang diminta
func HasAnyPermission(c *gin.Context, permissionCodes ...string) bool {
	// Superadmin bypass
	if strings.EqualFold(GetUserRole(c), "ADMIN") {
		return true
	}

	userPermissions := GetUserPermissions(c)
	permMap := make(map[string]struct{}, len(userPermissions))
	for _, p := range userPermissions {
		if p == "*" {
			return true
		}
		permMap[strings.ToLower(p)] = struct{}{}
	}

	for _, required := range permissionCodes {
		if _, ok := permMap[strings.ToLower(required)]; ok {
			return true
		}
	}
	return false
}

// GetUserID mengambil User ID dari context Gin
func GetUserID(c *gin.Context) string {
	if val, exists := c.Get(ContextUserID); exists {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserRole mengambil Role pengguna dari context Gin
func GetUserRole(c *gin.Context) string {
	if val, exists := c.Get(ContextUserRole); exists {
		if role, ok := val.(string); ok {
			return role
		}
	}
	return ""
}

// GetUsername mengambil Username dari context Gin
func GetUsername(c *gin.Context) string {
	if val, exists := c.Get(ContextUsername); exists {
		if username, ok := val.(string); ok {
			return username
		}
	}
	return ""
}

// GetUserPermissions mengambil slice permissions pengguna dari context Gin
func GetUserPermissions(c *gin.Context) []string {
	if val, exists := c.Get(ContextUserPermissions); exists {
		if perms, ok := val.([]string); ok {
			return perms
		}
		if permsInterface, ok := val.([]any); ok {
			result := make([]string, 0, len(permsInterface))
			for _, v := range permsInterface {
				if s, ok := v.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return []string{}
}
