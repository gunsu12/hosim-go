package departement

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
	departments := router.Group("/departments")
	{
		departments.POST("", h.Create)
		departments.GET("", h.List)
		departments.GET("/:id", h.GetByID)
		departments.PUT("/:id", h.Update)
		departments.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreateDepartementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.CreateDepartement(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidDepartementType) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrNameRequired) || errors.Is(err, ErrCodeRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Departemen berhasil dibuat", result)
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

	depts, count, err := h.service.ListDepartements(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar departemen berhasil diambil", gin.H{
		"data":  depts,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	dept, err := h.service.GetDepartementByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDepartementNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Departemen berhasil diambil", dept)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateDepartementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.UpdateDepartement(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrDepartementNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidDepartementType) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrNameRequired) || errors.Is(err, ErrCodeRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Departemen berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	err := h.service.DeleteDepartement(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrDepartementNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Departemen berhasil dihapus", nil)
}
