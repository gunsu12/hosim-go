package tariffcomponent

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
	components := router.Group("/tariff-components")
	{
		components.POST("", middleware.RequirePermission("tariff_component:create"), h.Create)
		components.GET("", middleware.RequirePermission("tariff_component:read"), h.List)
		components.GET("/:id", middleware.RequirePermission("tariff_component:read"), h.GetByID)
		components.PUT("/:id", middleware.RequirePermission("tariff_component:update"), h.Update)
		components.DELETE("/:id", middleware.RequirePermission("tariff_component:delete"), h.Delete)
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreateTariffComponentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.CreateTariffComponent(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrTariffComponentCodeAlreadyExists) {
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

	response.Success(ctx, http.StatusCreated, "Komponen tarif berhasil dibuat", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:   1,
		Limit:  10,
		Search:        ctx.Query("search"),
		ComponentType: ctx.Query("component_type"),
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

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	} else if params.Limit > 100 {
		params.Limit = 100
	}

	if ctx.Query("limit") != "" {
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "Parameter limit tidak valid", nil)
			return
		}
		params.Limit = limit
	}

	components, count, err := h.service.ListTariffComponents(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar komponen tarif berhasil diambil", gin.H{
		"data": components,
		"meta": gin.H{
			"total":       count,
			"page":        params.Page,
			"limit":       params.Limit,
			"total_pages": (count + int64(params.Limit) - 1) / int64(params.Limit),
		},
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	comp, err := h.service.GetTariffComponentByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTariffComponentNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Komponen tarif berhasil diambil", comp)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateTariffComponentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.UpdateTariffComponent(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrTariffComponentNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrTariffComponentCodeAlreadyExists) {
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

	response.Success(ctx, http.StatusOK, "Komponen tarif berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := getOperator(ctx)

	if err := h.service.DeleteTariffComponent(ctx.Request.Context(), id, operatorID); err != nil {
		if errors.Is(err, ErrTariffComponentNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Komponen tarif berhasil dihapus", nil)
}
