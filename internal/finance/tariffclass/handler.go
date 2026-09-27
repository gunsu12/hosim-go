package tariffclass

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
	return "SYSTEM"
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	classes := router.Group("/tariff-classes")
	{
		classes.POST("", middleware.RequirePermission("tariff_class:create"), h.Create)
		classes.GET("", middleware.RequirePermission("tariff_class:read"), h.List)
		classes.GET("/:id", middleware.RequirePermission("tariff_class:read"), h.GetByID)
		classes.PUT("/:id", middleware.RequirePermission("tariff_class:update"), h.Update)
		classes.DELETE("/:id", middleware.RequirePermission("tariff_class:delete"), h.Delete)
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreateTariffClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.CreateTariffClass(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrTariffClassCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Kelas tarif berhasil dibuat", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:   1,
		Limit:  10,
		Search: ctx.Query("search"),
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

	classes, count, err := h.service.ListTariffClasses(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar kelas tarif berhasil diambil", gin.H{
		"data":  classes,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	tc, err := h.service.GetTariffClassByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTariffClassNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Kelas tarif berhasil diambil", tc)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateTariffClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.UpdateTariffClass(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrTariffClassNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrTariffClassCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Kelas tarif berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := getOperator(ctx)

	err := h.service.DeleteTariffClass(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrTariffClassNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Kelas tarif berhasil dihapus", nil)
}
