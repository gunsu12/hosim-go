package room

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
	rooms := router.Group("/rooms")
	{
		rooms.POST("", h.Create)
		rooms.GET("", h.List)
		rooms.GET("/:id", h.GetByID)
		rooms.PUT("/:id", h.Update)
		rooms.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.CreateRoom(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrRoomCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidRoomType) || errors.Is(err, ErrNegativeCapacity) ||
			errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) ||
			errors.Is(err, ErrServiceUnitIDRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Ruangan berhasil dibuat", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:          1,
		Limit:         10,
		Search:        ctx.Query("search"),
		ServiceUnitID: ctx.Query("service_unit_id"),
		RoomType:      ctx.Query("room_type"),
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

	rooms, count, err := h.service.ListRooms(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar ruangan berhasil diambil", gin.H{
		"data":  rooms,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	room, err := h.service.GetRoomByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRoomNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Ruangan berhasil diambil", room)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	result, err := h.service.UpdateRoom(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrRoomNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrRoomCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidRoomType) || errors.Is(err, ErrNegativeCapacity) ||
			errors.Is(err, ErrCodeRequired) || errors.Is(err, ErrNameRequired) ||
			errors.Is(err, ErrServiceUnitIDRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Ruangan berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := ctx.GetHeader("X-User-ID")
	if operatorID == "" {
		operatorID = "SYSTEM"
	}

	err := h.service.DeleteRoom(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrRoomNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Ruangan berhasil dihapus", nil)
}
