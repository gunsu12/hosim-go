package journal

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hosim-go/internal/accounting/fiscalperiod"
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
	journals := router.Group("/accounting/journals")
	{
		journals.POST("", middleware.RequirePermission("accounting:journal:create"), h.Create)
		journals.GET("", middleware.RequirePermission("accounting:journal:read"), h.List)
		journals.GET("/:id", middleware.RequirePermission("accounting:journal:read"), h.GetByID)
		journals.PUT("/:id", middleware.RequirePermission("accounting:journal:update"), h.Update)
		journals.DELETE("/:id", middleware.RequirePermission("accounting:journal:delete"), h.Delete)
		journals.POST("/:id/post", middleware.RequirePermission("accounting:journal:post"), h.Post)
		journals.POST("/:id/reverse", middleware.RequirePermission("accounting:journal:reverse"), h.Reverse)
	}
}

// Create menangani pembuatan draf voucher jurnal manual baru.
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateJournalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(ctx)
	res, err := h.service.CreateDraft(ctx.Request.Context(), req, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, "Draf jurnal berhasil dibuat", res)
}

// List mengambil daftar jurnal terpaginasi dengan filter.
func (h *Handler) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	params := ListJournalParams{
		Page:   page,
		Limit:  limit,
		Search: strings.TrimSpace(ctx.Query("search")),
	}

	if s := strings.TrimSpace(ctx.Query("start_date")); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			params.StartDate = &t
		}
	}
	if e := strings.TrimSpace(ctx.Query("end_date")); e != "" {
		if t, err := time.Parse("2006-01-02", e); err == nil {
			params.EndDate = &t
		}
	}
	if fp := strings.TrimSpace(ctx.Query("fiscal_period")); fp != "" {
		params.FiscalPeriod = &fp
	}
	if st := strings.TrimSpace(ctx.Query("status")); st != "" {
		status := JournalStatus(strings.ToUpper(st))
		if status.IsValid() {
			params.Status = &status
		}
	}
	if et := strings.TrimSpace(ctx.Query("entry_type")); et != "" {
		eType := JournalType(strings.ToUpper(et))
		if eType.IsValid() {
			params.EntryType = &eType
		}
	}

	res, total, err := h.service.ListJournals(ctx.Request.Context(), params)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.SuccessWithMeta(ctx, http.StatusOK, "Daftar jurnal berhasil dimuat", res, map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetByID mengambil detail satu voucher jurnal beserta baris rinciannya.
func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := h.service.GetJournal(ctx.Request.Context(), id)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Detail jurnal berhasil dimuat", res)
}

// Update memperbarui isi draf voucher jurnal.
func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateJournalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(ctx)
	res, err := h.service.UpdateDraft(ctx.Request.Context(), id, req, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Draf jurnal berhasil diperbarui", res)
}

// Delete menghapus draf jurnal (hanya boleh jika status DRAFT).
func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operator := getOperator(ctx)

	if err := h.service.DeleteDraft(ctx.Request.Context(), id, operator); err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Draf jurnal berhasil dihapus", nil)
}

// Post mengesahkan draf jurnal menjadi POSTED ke Buku Besar.
func (h *Handler) Post(ctx *gin.Context) {
	id := ctx.Param("id")
	operator := getOperator(ctx)

	res, err := h.service.PostJournal(ctx.Request.Context(), id, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Jurnal berhasil disahkan (POSTED) ke Buku Besar", res)
}

// Reverse membatalkan jurnal berstatus POSTED dengan menerbitkan jurnal pembalik.
func (h *Handler) Reverse(ctx *gin.Context) {
	id := ctx.Param("id")
	var req ReverseJournalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: alasan pembatalan wajib diisi", err.Error())
		return
	}

	operator := getOperator(ctx)
	res, err := h.service.ReverseJournal(ctx.Request.Context(), id, req, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Jurnal berhasil di-reverse (dibatalkan dengan jurnal pembalik)", res)
}

func (h *Handler) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrJournalNotFound):
		response.Error(ctx, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, ErrJournalUnbalanced),
		errors.Is(err, ErrJournalMinimumLines),
		errors.Is(err, ErrJournalNotDraft),
		errors.Is(err, ErrJournalAlreadyPosted),
		errors.Is(err, ErrJournalAlreadyReversed),
		errors.Is(err, ErrCannotReverseUnpostedJournal),
		errors.Is(err, ErrInvalidEntryPosition),
		errors.Is(err, ErrInvalidAmount),
		errors.Is(err, ErrReversalReasonRequired),
		errors.Is(err, ErrAccountNotPostable),
		errors.Is(err, ErrBackdateReferenceRequired),
		errors.Is(err, fiscalperiod.ErrPeriodClosed),
		errors.Is(err, fiscalperiod.ErrPeriodLocked):
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, ErrTrialBalanceCorrupted):
		response.Error(ctx, http.StatusUnprocessableEntity, err.Error(), nil)
	default:
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
}
