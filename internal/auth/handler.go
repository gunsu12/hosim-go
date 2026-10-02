package auth

import (
	"errors"
	"net/http"
	"strconv"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   Service
	jwtSecret string
}

func NewHandler(service Service, jwtSecret string) *Handler {
	return &Handler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

func getOperator(c *gin.Context) string {
	if op := c.GetHeader("X-User-ID"); op != "" {
		return op
	}
	if username := middleware.GetUsername(c); username != "" {
		return username
	}
	if uid := middleware.GetUserID(c); uid != "" {
		return uid
	}
	return "SYSTEM"
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.RefreshToken)
		authGroup.POST("/logout", h.Logout)

		// Endpoint terproteksi yang membutuhkan token JWT
		protected := authGroup.Group("")
		protected.Use(middleware.AuthMiddleware(h.jwtSecret))
		{
			protected.GET("/me", h.GetMe)

			// Manajemen Pengguna (Users)
			usersGroup := protected.Group("/users")
			{
				usersGroup.GET("", middleware.RequirePermission("user:read"), h.ListUsers)
				usersGroup.POST("", middleware.RequirePermission("user:create"), h.CreateUser)
				usersGroup.GET("/:id", middleware.RequirePermission("user:read"), h.GetUserByID)
				usersGroup.PUT("/:id", middleware.RequirePermission("user:update"), h.UpdateUser)
				usersGroup.DELETE("/:id", middleware.RequirePermission("user:delete"), h.DeleteUser)
			}

			// Manajemen Peran (Roles)
			rolesGroup := protected.Group("/roles")
			{
				rolesGroup.GET("", middleware.RequirePermission("role:read"), h.GetRoles)
				rolesGroup.POST("", middleware.RequirePermission("role:create"), h.CreateRole)
				rolesGroup.GET("/:id", middleware.RequirePermission("role:read"), h.GetRoleByID)
				rolesGroup.PUT("/:id", middleware.RequirePermission("role:update"), h.UpdateRole)
				rolesGroup.DELETE("/:id", middleware.RequirePermission("role:delete"), h.DeleteRole)
			}

			// Manajemen Hak Akses (Permissions)
			permsGroup := protected.Group("/permissions")
			{
				permsGroup.GET("", middleware.RequirePermission("role:read"), h.GetPermissions)
				permsGroup.POST("", middleware.RequireRole("ADMIN"), h.CreatePermission)
				permsGroup.PUT("/:id", middleware.RequireRole("ADMIN"), h.UpdatePermission)
				permsGroup.DELETE("/:id", middleware.RequireRole("ADMIN"), h.DeletePermission)
			}
		}
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	result, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUsernameTaken) || errors.Is(err, ErrEmailTaken) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Pendaftaran pengguna berhasil", result)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Username dan password wajib diisi", err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrUserInactive) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", result)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Refresh token wajib dikirim", err.Error())
		return
	}

	result, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			response.Error(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrUserInactive) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Token berhasil diperbarui", result)
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal memproses logout: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Logout berhasil", nil)
}

func (h *Handler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Pengguna belum terautentikasi", nil)
		return
	}

	result, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Profil pengguna berhasil diambil", result)
}

// ========================================================
// User Handlers
// ========================================================

func (h *Handler) ListUsers(c *gin.Context) {
	search := c.Query("search")
	roleID := c.Query("role_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		b := activeStr == "true" || activeStr == "1"
		isActive = &b
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), search, roleID, isActive, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar pengguna: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Daftar pengguna berhasil diambil", gin.H{
		"data":  users,
		"count": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	result, err := h.service.CreateUser(c.Request.Context(), req, operator)
	if err != nil {
		if errors.Is(err, ErrUsernameTaken) || errors.Is(err, ErrEmailTaken) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal membuat pengguna: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Pengguna berhasil dibuat", result)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	result, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Detail pengguna berhasil diambil", result)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	result, err := h.service.UpdateUser(c.Request.Context(), id, req, operator)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrEmailTaken) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal memperbarui pengguna: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Data pengguna berhasil diperbarui", result)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	operator := getOperator(c)

	if err := h.service.DeleteUser(c.Request.Context(), id, operator); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal menghapus pengguna: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Pengguna berhasil dinonaktifkan/dihapus", nil)
}

// ========================================================
// Role Handlers
// ========================================================

func (h *Handler) GetRoles(c *gin.Context) {
	roles, err := h.service.GetRoles(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data role: "+err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Daftar role berhasil diambil", roles)
}

func (h *Handler) GetRoleByID(c *gin.Context) {
	id := c.Param("id")
	role, err := h.service.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Detail role berhasil diambil", role)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	role, err := h.service.CreateRole(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrRoleCodeTaken) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal membuat role: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Role berhasil dibuat", role)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	role, err := h.service.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal memperbarui role: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Role berhasil diperbarui", role)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCannotDeleteAdmin) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Role berhasil dihapus", nil)
}

// ========================================================
// Permission Handlers
// ========================================================

func (h *Handler) GetPermissions(c *gin.Context) {
	perms, err := h.service.GetPermissions(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data permissions: "+err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Daftar permissions berhasil diambil", perms)
}

func (h *Handler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	perm, err := h.service.CreatePermission(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrPermCodeTaken) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal membuat permission: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Permission berhasil dibuat", perm)
}

func (h *Handler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	perm, err := h.service.UpdatePermission(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal memperbarui permission: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Permission berhasil diperbarui", perm)
}

func (h *Handler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePermission(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal menghapus permission: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Permission berhasil dihapus", nil)
}
