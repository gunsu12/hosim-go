package account

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/enums"
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
	accounts := router.Group("/accounting/accounts")
	{
		accounts.POST("", middleware.RequirePermission("accounting:create"), h.Create)
		accounts.GET("", middleware.RequirePermission("accounting:read"), h.List)
		accounts.GET("/tree", middleware.RequirePermission("accounting:read"), h.GetTree)
		accounts.GET("/postable", middleware.RequirePermission("accounting:read"), h.ListPostable)
		accounts.GET("/treasury", middleware.RequirePermission("accounting:read"), h.ListTreasury)
		accounts.GET("/:id", middleware.RequirePermission("accounting:read"), h.GetByID)
		accounts.PUT("/:id", middleware.RequirePermission("accounting:update"), h.Update)
		accounts.DELETE("/:id", middleware.RequirePermission("accounting:delete"), h.Delete)
	}
}

// Create menangani pembuatan akun baru
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(ctx)
	result, err := h.service.CreateAccount(ctx.Request.Context(), req, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, "Akun bagan berhasil dibuat", result)
}

// List menangani query daftar akun datar dengan filter & paginasi
func (h *Handler) List(ctx *gin.Context) {
	params := AccountListParams{
		Page:   1,
		Limit:  10,
		Search: strings.TrimSpace(ctx.Query("search")),
	}

	if ctx.Query("page") != "" {
		if page, err := strconv.Atoi(ctx.Query("page")); err == nil && page > 0 {
			params.Page = page
		}
	}

	if ctx.Query("limit") != "" {
		if limit, err := strconv.Atoi(ctx.Query("limit")); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100
			}
			params.Limit = limit
		}
	}

	if t := strings.TrimSpace(ctx.Query("type")); t != "" {
		accType := enums.AccountType(strings.ToUpper(t))
		if IsValidType(accType) {
			params.Type = &accType
		}
	}

	if p := strings.TrimSpace(ctx.Query("position")); p != "" {
		pos := enums.AccountPosition(strings.ToUpper(p))
		if IsValidPosition(pos) {
			params.Position = &pos
		}
	}

	if parentID := strings.TrimSpace(ctx.Query("parent_id")); parentID != "" {
		params.ParentID = &parentID
	}

	if postable := ctx.Query("is_postable"); postable != "" {
		val := postable == "true" || postable == "1"
		params.IsPostable = &val
	}

	if treasury := ctx.Query("is_treasury"); treasury != "" {
		val := treasury == "true" || treasury == "1"
		params.IsTreasury = &val
	}

	if active := ctx.Query("is_active"); active != "" {
		val := active == "true" || active == "1"
		params.IsActive = &val
	}

	if lvlStr := ctx.Query("level"); lvlStr != "" {
		if lvl, err := strconv.Atoi(lvlStr); err == nil && lvl > 0 {
			params.Level = &lvl
		}
	}

	accounts, total, err := h.service.List(ctx.Request.Context(), params)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))
	if total == 0 {
		totalPages = 0
	}

	meta := response.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	response.SuccessWithMeta(ctx, http.StatusOK, "Daftar akun berhasil diambil", accounts, meta)
}

// GetTree menangani pengambilan struktur pohon akun hierarkis
func (h *Handler) GetTree(ctx *gin.Context) {
	activeOnly := ctx.Query("active_only") == "true" || ctx.Query("active_only") == "1"

	tree, err := h.service.GetTree(ctx.Request.Context(), activeOnly)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Pohon bagan akun berhasil diambil", tree)
}

// ListPostable mengembalikan akun daun yang siap digunakan untuk penjurnalan dan pemetaan transaksi
func (h *Handler) ListPostable(ctx *gin.Context) {
	accounts, err := h.service.ListPostable(ctx.Request.Context())
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar akun postable berhasil diambil", accounts)
}

// ListTreasury mengembalikan akun kas dan bank perbendaharaan aktif
func (h *Handler) ListTreasury(ctx *gin.Context) {
	accounts, err := h.service.ListTreasury(ctx.Request.Context())
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Daftar akun perbendaharaan kas/bank berhasil diambil", accounts)
}

// GetByID mengambil detail satu akun berdasarkan UUID
func (h *Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	acc, err := h.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Detail akun berhasil diambil", acc)
}

// Update mengubah atribut metadata akun
func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(ctx)
	result, err := h.service.UpdateAccount(ctx.Request.Context(), id, req, operator)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Akun bagan berhasil diperbarui", result)
}

// Delete melakukan soft delete akun dengan validasi tidak ada anak dan tidak terikat tarif
func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	operator := getOperator(ctx)

	if err := h.service.DeleteAccount(ctx.Request.Context(), id, operator); err != nil {
		h.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Akun bagan berhasil dihapus", nil)
}

func (h *Handler) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrAccountNotFound):
		response.Error(ctx, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, ErrAccountCodeAlreadyExists):
		response.Error(ctx, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, ErrAccountCodeRequired),
		errors.Is(err, ErrAccountNameRequired),
		errors.Is(err, ErrInvalidAccountLevel),
		errors.Is(err, ErrParentAccountNotFound),
		errors.Is(err, ErrInvalidAccountType),
		errors.Is(err, ErrInvalidAccountPosition),
		errors.Is(err, ErrTypeMismatchWithParent),
		errors.Is(err, ErrParentCannotBeSelf),
		errors.Is(err, ErrTreasuryMustBeAsset),
		errors.Is(err, ErrTreasuryMustBePostable),
		errors.Is(err, ErrNonPostableAccount),
		errors.Is(err, ErrAccountInactive):
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, ErrCannotDeleteAccountWithChildren),
		errors.Is(err, ErrCannotDeleteAccountInUse):
		response.Error(ctx, http.StatusUnprocessableEntity, err.Error(), nil)
	default:
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
}
