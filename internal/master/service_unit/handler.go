package serviceunit

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
	units := router.Group("/service-units")
	{
		units.POST("", h.Create)
		units.GET("", h.List)
		units.GET("/:id", h.GetByID)
		units.PUT("/:id", h.Update)
		units.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreateServiceUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.CreateServiceUnit(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrServiceUnitCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Unit layanan berhasil dibuat", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:          1,
		Limit:         10,
		Search:        ctx.Query("search"),
		DepartementID: ctx.Query("departement_id"),
	}

	if ctx.Query("is_registration_target") != "" {
		val, err := strconv.ParseBool(ctx.Query("is_registration_target"))
		if err == nil {
			params.IsRegistrationTarget = &val
		}
	}

	if ctx.Query("is_active") != "" {
		val, err := strconv.ParseBool(ctx.Query("is_active"))
		if err == nil {
			params.IsActive = &val
		}
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

	units, count, err := h.service.ListServiceUnits(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar unit layanan berhasil diambil", gin.H{
		"data":  units,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	su, err := h.service.GetServiceUnitByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrServiceUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Unit layanan berhasil diambil", su)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateServiceUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.UpdateServiceUnit(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrServiceUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrServiceUnitCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Unit layanan berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	err := h.service.DeleteServiceUnit(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrServiceUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Unit layanan berhasil dihapus", nil)
}
