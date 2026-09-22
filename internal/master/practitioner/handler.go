package practitioner

import (
	"errors"
	"net/http"
	"strconv"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
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
	return "petugas-sdm"
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	practitioners := router.Group("/practitioners")
	{
		// Cipta, Ubah & Hapus data dokter dibatasi ketat
		practitioners.POST("", middleware.RequirePermission("practitioner:create"), h.Create)
		practitioners.PUT("/:id", middleware.RequirePermission("practitioner:update"), h.Update)
		practitioners.DELETE("/:id", middleware.RequirePermission("practitioner:delete"), h.Delete)
		practitioners.GET("/nik/:nik", middleware.RequirePermission("practitioner:read"), h.GetByNIK)
		practitioners.GET("/nip/:nip", middleware.RequirePermission("practitioner:read"), h.GetByNIP)

		// Opsi 1: GET /practitioners diizinkan untuk Admin Master ATAU petugas yang butuh data dokter
		practitioners.GET("", middleware.RequireAnyPermission(
			"practitioner:read",   // Admin Master Data
			"outpatient:register", // Suster/Paramedis Pendaftaran Rawat Jalan
			"appointment:view",    // Petugas Janji Temu
		), h.List)

		practitioners.GET("/:id", middleware.RequireAnyPermission(
			"practitioner:read",
			"outpatient:register",
			"appointment:view",
		), h.GetByID)
	}

	router.GET("/professions", middleware.RequireAnyPermission("practitioner:read", "practitioner:create", "outpatient:register"), h.ListProfessions)
	router.GET("/specialties", middleware.RequireAnyPermission("practitioner:read", "practitioner:create", "outpatient:register"), h.ListSpecialties)
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePractitionerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.CreatePractitioner(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrNIKAlreadyExists) || errors.Is(err, ErrNIPAlreadyExists) || errors.Is(err, ErrIhsAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidNIKFormat) || errors.Is(err, ErrInvalidGender) ||
			errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrInvalidSIPExpiryDateFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Tenaga medis berhasil didaftarkan", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:   1,
		Limit:  10,
		Search: ctx.Query("search"),
	}
	if ctx.Query("page") != "" {
		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "Parameter page tidak valid", nil)
			return
		}
		params.Page = page
	}
	if ctx.Query("limit") != "" {
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "Parameter limit tidak valid", nil)
			return
		}
		params.Limit = limit
	}

	practitioners, count, err := h.service.ListPractitioners(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar tenaga medis berhasil diambil", gin.H{
		"data":  practitioners,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := h.service.GetPractitionerByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPractitionerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Data tenaga medis ditemukan", result)
}

func (h *Handler) GetByNIK(ctx *gin.Context) {
	nik := ctx.Param("nik")
	result, err := h.service.GetPractitionerByNIK(ctx.Request.Context(), nik)
	if err != nil {
		if errors.Is(err, ErrPractitionerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Data tenaga medis ditemukan", result)
}

func (h *Handler) GetByNIP(ctx *gin.Context) {
	nip := ctx.Param("nip")
	result, err := h.service.GetPractitionerByNIP(ctx.Request.Context(), nip)
	if err != nil {
		if errors.Is(err, ErrPractitionerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Data tenaga medis ditemukan", result)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdatePractitionerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.UpdatePractitioner(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrPractitionerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrNIKAlreadyExists) || errors.Is(err, ErrNIPAlreadyExists) || errors.Is(err, ErrIhsAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidNIKFormat) || errors.Is(err, ErrInvalidGender) ||
			errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrInvalidSIPExpiryDateFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Data tenaga medis berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := getOperator(ctx)

	err := h.service.DeletePractitioner(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrPractitionerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Data tenaga medis berhasil dihapus", nil)
}

func (h *Handler) ListProfessions(ctx *gin.Context) {
	professions, err := h.service.ListProfessions(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "Daftar profesi berhasil diambil", professions)
}

func (h *Handler) ListSpecialties(ctx *gin.Context) {
	specialties, err := h.service.ListSpecialties(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "Daftar spesialisasi berhasil diambil", specialties)
}

