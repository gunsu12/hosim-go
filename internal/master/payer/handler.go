package payer

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
	payers := router.Group("/payers")
	{
		payers.POST("", h.Create)
		payers.GET("", h.List)
		payers.GET("/:id", h.GetByID)
		payers.PUT("/:id", h.Update)
		payers.DELETE("/:id", h.Delete)
	}

	router.GET("/payer-types", h.ListPayerTypes)
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.CreatePayer(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrPayerCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPayerTypeNotFound) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrNameRequired) || errors.Is(err, ErrCodeRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Penjamin berhasil dibuat", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:        1,
		Limit:       10,
		Search:      ctx.Query("search"),
		PayerTypeID: ctx.Query("payer_type_id"),
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

	payers, count, err := h.service.ListPayers(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar penjamin berhasil diambil", gin.H{
		"data":  payers,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	payer, err := h.service.GetPayerByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPayerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Penjamin berhasil diambil", payer)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdatePayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.UpdatePayer(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrPayerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPayerCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPayerTypeNotFound) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrNameRequired) || errors.Is(err, ErrCodeRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Penjamin berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	err := h.service.DeletePayer(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrPayerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPayerImmutable) {
			response.Error(ctx, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Penjamin berhasil dihapus", nil)
}

func (h *Handler) ListPayerTypes(ctx *gin.Context) {
	types, err := h.service.ListPayerTypes(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar tipe penjamin berhasil diambil", types)
}
