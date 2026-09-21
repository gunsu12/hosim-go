package practitioner

import (
	"errors"
	"net/http"
	"strconv"

	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	practitioners := router.Group("/practitioners")
	{
		practitioners.POST("", h.Create)
		practitioners.GET("", h.List)
		practitioners.GET("/:id", h.GetByID)
		practitioners.GET("/nik/:nik", h.GetByNIK)
		practitioners.GET("/nip/:nip", h.GetByNIP)
		practitioners.PUT("/:id", h.Update)
		practitioners.DELETE("/:id", h.Delete)
	}

	router.GET("/professions", h.ListProfessions)
	router.GET("/specialties", h.ListSpecialties)
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePractitionerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "petugas-sdm"
	}

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

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "petugas-sdm"
	}

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
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "petugas-sdm"
	}

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

