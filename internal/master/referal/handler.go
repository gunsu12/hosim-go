package referal

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
	referals := router.Group("/referals")
	{
		referals.POST("", h.Create)
		referals.GET("", h.List)
		referals.GET("/:id", h.GetByID)
		referals.PUT("/:id", h.Update)
		referals.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreateReferalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.CreateReferal(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrReferalCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidReferalType) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Rujukan berhasil dibuat", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:   1,
		Limit:  10,
		Search: ctx.Query("search"),
		Type:   ctx.Query("type"),
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

	refs, count, err := h.service.ListReferals(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar rujukan berhasil diambil", gin.H{
		"data":  refs,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	ref, err := h.service.GetReferalByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrReferalNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Rujukan berhasil diambil", ref)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateReferalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.UpdateReferal(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrReferalNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrReferalCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidReferalType) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Rujukan berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	err := h.service.DeleteReferal(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrReferalNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Rujukan berhasil dihapus", nil)
}
