package auth

import (
	"errors"
	"net/http"

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
			protected.GET("/roles", h.GetRoles)
			protected.GET("/permissions", h.GetPermissions)
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
	// Logout tetap bisa dipanggil walau tanpa body json, tapi jika ada refresh_token kita hapus dari database
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

func (h *Handler) GetRoles(c *gin.Context) {
	roles, err := h.service.GetRoles(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data role", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Daftar role berhasil diambil", roles)
}

func (h *Handler) GetPermissions(c *gin.Context) {
	perms, err := h.service.GetPermissions(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data permissions", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Daftar permissions berhasil diambil", perms)
}
