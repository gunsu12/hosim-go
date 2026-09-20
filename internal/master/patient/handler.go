package patient

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
	patients := router.Group("/patients")
	{
		patients.POST("", h.Create)
		patients.GET("", h.List)
		patients.GET("/:id", h.GetByID)
		patients.GET("/nik/:nik", h.GetByNIK)
		patients.GET("/rm/:rmNo", h.GetByMedicalRecordNo)
		patients.PUT("/:id", h.Update)
		patients.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePatientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "petugas-loket"
	}

	patient, err := h.service.RegisterPatient(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrNIKAlreadyExists) || errors.Is(err, ErrMedicalRecordNoExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidDateFormat) || errors.Is(err, ErrBirthDateInFuture) ||
			errors.Is(err, ErrInvalidNIKFormat) || errors.Is(err, ErrInvalidFamilyCardFormat) ||
			errors.Is(err, ErrInvalidGender) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrInvalidDeceasedDate) || errors.Is(err, ErrMaxMedicalRecordExceeded) ||
			errors.Is(err, ErrInvalidMedicalRecordFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusCreated, "Pasien berhasil didaftarkan", patient)

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
			response.Error(ctx, http.StatusBadRequest, "Invalid page number", nil)
			return
		}
		params.Page = page
	}
	if ctx.Query("limit") != "" {
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "Invalid limit", nil)
			return
		}
		params.Limit = limit
	}

	patients, count, err := h.service.ListPatients(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "List pasien berhasil diambil", gin.H{
		"data":  patients,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	patient, err := h.service.GetPatientByID(ctx.Request.Context(), id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Pasien tidak ditemukan", nil)
		return
	}
	response.Success(ctx, http.StatusOK, "Pasien berhasil diambil", patient)
}

func (h *Handler) GetByNIK(ctx *gin.Context) {
	nik := ctx.Param("nik")
	patient, err := h.service.GetPatientByNIK(ctx.Request.Context(), nik)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Pasien tidak ditemukan", nil)
		return
	}
	response.Success(ctx, http.StatusOK, "Pasien berhasil diambil", patient)
}

func (h *Handler) GetByMedicalRecordNo(ctx *gin.Context) {
	rmNo := ctx.Param("rmNo")
	patient, err := h.service.GetPatientByMedicalRecordNo(ctx.Request.Context(), rmNo)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Pasien tidak ditemukan", nil)
		return
	}
	response.Success(ctx, http.StatusOK, "Pasien berhasil diambil", patient)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdatePatientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	// Ambil ID operator dari header/token (sementara default: "petugas-loket")
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "petugas-loket"
	}

	patient, err := h.service.UpdatePatient(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrPatientNotFound) {
			response.Error(ctx, http.StatusNotFound, "Pasien tidak ditemukan", err.Error())
			return
		}
		if errors.Is(err, ErrNIKAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidDateFormat) || errors.Is(err, ErrBirthDateInFuture) ||
			errors.Is(err, ErrInvalidNIKFormat) || errors.Is(err, ErrInvalidFamilyCardFormat) ||
			errors.Is(err, ErrInvalidGender) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrInvalidDeceasedDate) || errors.Is(err, ErrMaxMedicalRecordExceeded) ||
			errors.Is(err, ErrInvalidMedicalRecordFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "Pasien berhasil diupdate", patient)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	// Ambil ID operator dari header/token (sementara default: "petugas-loket")
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "petugas-loket"
	}

	err := h.service.DeletePatient(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrPatientNotFound) {
			response.Error(ctx, http.StatusNotFound, "Pasien tidak ditemukan", err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "Pasien berhasil dihapus", nil)
}
