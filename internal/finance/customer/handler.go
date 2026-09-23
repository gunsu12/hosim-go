package customer

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
	customers := router.Group("/customers")
	{
		customers.POST("", middleware.RequirePermission("customer:create"), h.Create)
		customers.GET("", middleware.RequirePermission("customer:read"), h.List)
		customers.GET("/:id", middleware.RequirePermission("customer:read"), h.GetByID)
		customers.PUT("/:id", middleware.RequirePermission("customer:update"), h.Update)
		customers.DELETE("/:id", middleware.RequirePermission("customer:delete"), h.Delete)
	}

	router.GET("/customer-types", middleware.RequirePermission("customer:read"), h.ListCustomerTypes)
}

func (h *Handler) Create(ctx *gin.Context) {
	var req CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.CreateCustomer(ctx.Request.Context(), req, operatorID)
	if err != nil {
		if errors.Is(err, ErrCustomerCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCustomerTypeNotFound) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrNameRequired) || errors.Is(err, ErrCodeRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "Customer berhasil dibuat", result)
}

func (h *Handler) List(ctx *gin.Context) {
	params := ListParams{
		Page:           1,
		Limit:          10,
		Search:         ctx.Query("search"),
		CustomerTypeID: ctx.Query("customer_type_id"),
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

	customers, count, err := h.service.ListCustomers(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar customer berhasil diambil", gin.H{
		"data":  customers,
		"count": count,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	customer, err := h.service.GetCustomerByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Customer berhasil diambil", customer)
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operatorID := getOperator(ctx)

	result, err := h.service.UpdateCustomer(ctx.Request.Context(), id, req, operatorID)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCustomerCodeAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCustomerTypeNotFound) || errors.Is(err, ErrInvalidEmail) ||
			errors.Is(err, ErrNameRequired) || errors.Is(err, ErrCodeRequired) {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Customer berhasil diperbarui", result)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operatorID := getOperator(ctx)

	err := h.service.DeleteCustomer(ctx.Request.Context(), id, operatorID)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			response.Error(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCustomerImmutable) {
			response.Error(ctx, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Customer berhasil dihapus", nil)
}

func (h *Handler) ListCustomerTypes(ctx *gin.Context) {
	types, err := h.service.ListCustomerTypes(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar tipe customer berhasil diambil", types)
}
