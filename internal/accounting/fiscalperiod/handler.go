package fiscalperiod

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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
	if username := middleware.GetUsername(c); username != "" {
		return username
	}
	if uid := middleware.GetUserID(c); uid != "" {
		return uid
	}
	return "SYSTEM"
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// 1. Endpoint Tahun Fiskal
	years := router.Group("/accounting/fiscal-years")
	{
		years.POST("", middleware.RequirePermission("accounting:period:manage"), h.CreateYear)
		years.GET("", middleware.RequirePermission("accounting:period:read"), h.ListYears)
		years.GET("/:id", middleware.RequirePermission("accounting:period:read"), h.GetYearByID)
		years.PUT("/:id/close", middleware.RequirePermission("accounting:period:close"), h.CloseYear)
	}

	// 2. Endpoint Periode Bulanan
	periods := router.Group("/accounting/fiscal-periods")
	{
		periods.GET("", middleware.RequirePermission("accounting:period:read"), h.ListPeriods)
		periods.GET("/current", middleware.RequirePermission("accounting:period:read"), h.GetCurrentPeriod)
		periods.GET("/:id", middleware.RequirePermission("accounting:period:read"), h.GetPeriodByID)
		periods.PUT("/:id/lock", middleware.RequirePermission("accounting:period:manage"), h.LockPeriod)
		periods.PUT("/:id/unlock", middleware.RequirePermission("accounting:period:manage"), h.UnlockPeriod)
		periods.PUT("/:id/close", middleware.RequirePermission("accounting:period:close"), h.ClosePeriod)
	}
}

// CreateYear membuat tahun buku baru beserta 12 periode bulanannya.
func (h *Handler) CreateYear(ctx *gin.Context) {
	var req CreateFiscalYearRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(ctx)
	res, err := h.service.CreateFiscalYear(ctx.Request.Context(), req, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, "Tahun fiskal dan 12 periode bulanan berhasil dibuat", res)
}

// ListYears mengambil daftar tahun buku yang terdaftar.
func (h *Handler) ListYears(ctx *gin.Context) {
	res, err := h.service.ListFiscalYears(ctx.Request.Context())
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar tahun fiskal berhasil dimuat", res)
}

// GetYearByID mengambil detail satu tahun fiskal beserta daftar periodenya.
func (h *Handler) GetYearByID(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := h.service.GetFiscalYear(ctx.Request.Context(), id)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Detail tahun fiskal berhasil dimuat", res)
}

// CloseYear menutup buku tahunan (seluruh 12 bulan harus sudah closed).
func (h *Handler) CloseYear(ctx *gin.Context) {
	id := ctx.Param("id")
	operator := getOperator(ctx)

	res, err := h.service.CloseFiscalYear(ctx.Request.Context(), id, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Tahun fiskal berhasil ditutup buku permanen", res)
}

// ListPeriods mengambil daftar periode bulanan dengan filter tahun dan status.
func (h *Handler) ListPeriods(ctx *gin.Context) {
	var yearPtr *int
	if yStr := ctx.Query("year"); yStr != "" {
		if y, err := strconv.Atoi(yStr); err == nil {
			yearPtr = &y
		}
	}

	var statusPtr *string
	if s := strings.TrimSpace(ctx.Query("status")); s != "" {
		upper := strings.ToUpper(s)
		statusPtr = &upper
	}

	res, err := h.service.ListFiscalPeriods(ctx.Request.Context(), yearPtr, statusPtr)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar periode fiskal berhasil dimuat", res)
}

// GetCurrentPeriod mengambil periode fiskal yang aktif untuk hari ini.
func (h *Handler) GetCurrentPeriod(ctx *gin.Context) {
	res, err := h.service.GetCurrentPeriod(ctx.Request.Context())
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Periode fiskal aktif berhasil dimuat", res)
}

// GetPeriodByID mengambil detail satu periode bulanan.
func (h *Handler) GetPeriodByID(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := h.service.GetFiscalPeriod(ctx.Request.Context(), id)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Detail periode fiskal berhasil dimuat", res)
}

// LockPeriod mengunci sementara periode bulanan (LOCKED).
func (h *Handler) LockPeriod(ctx *gin.Context) {
	id := ctx.Param("id")
	var req FiscalPeriodActionRequest
	_ = ctx.ShouldBindJSON(&req)

	operator := getOperator(ctx)
	res, err := h.service.LockFiscalPeriod(ctx.Request.Context(), id, req.Notes, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Periode fiskal berhasil dikunci sementara (LOCKED)", res)
}

// UnlockPeriod membuka kembali periode yang berstatus LOCKED menjadi OPEN.
func (h *Handler) UnlockPeriod(ctx *gin.Context) {
	id := ctx.Param("id")
	var req FiscalPeriodActionRequest
	_ = ctx.ShouldBindJSON(&req)

	operator := getOperator(ctx)
	res, err := h.service.UnlockFiscalPeriod(ctx.Request.Context(), id, req.Notes, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Kunci periode fiskal berhasil dibuka (OPEN)", res)
}

// ClosePeriod menutup buku permanen sebuah periode bulanan (CLOSED).
func (h *Handler) ClosePeriod(ctx *gin.Context) {
	id := ctx.Param("id")
	var req FiscalPeriodActionRequest
	_ = ctx.ShouldBindJSON(&req)

	operator := getOperator(ctx)
	res, err := h.service.CloseFiscalPeriod(ctx.Request.Context(), id, req.Notes, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Periode fiskal berhasil ditutup buku permanen (CLOSED)", res)
}

func (h *Handler) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrYearNotFound), errors.Is(err, ErrPeriodNotFound):
		response.Error(ctx, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, ErrYearAlreadyExists):
		response.Error(ctx, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, ErrInvalidYear),
		errors.Is(err, ErrPeriodNotLocked),
		errors.Is(err, ErrPeriodAlreadyClosed),
		errors.Is(err, ErrCannotLockClosedPeriod),
		errors.Is(err, ErrYearAlreadyClosed):
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, ErrCannotCloseYearWithOpenPeriods):
		response.Error(ctx, http.StatusUnprocessableEntity, err.Error(), nil)
	default:
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
}
