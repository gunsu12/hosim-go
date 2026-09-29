package item

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"
	"uuid"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// ==========================================
// Product Line DTOs
// ==========================================

type CreateProductLineRequest struct {
	Code                            string  `json:"code" binding:"required"`
	Name                            string  `json:"name" binding:"required"`
	InventoryCOACode                *string `json:"inventory_coa_code"`
	CogsCOACode                     *string `json:"cogs_coa_code"`
	PurchaseDiscountCOACode         *string `json:"purchase_discount_coa_code"`
	PurchaseTaxCOACode              *string `json:"purchase_tax_coa_code"`
	AssetCOACode                    *string `json:"asset_coa_code"`
	AssetAccumulationCOACode        *string `json:"asset_accumulation_coa_code"`
	AssetDepreciationExpenseCOACode *string `json:"asset_depreciation_expense_coa_code"`
	IsActive                        *bool   `json:"is_active"`
}

type UpdateProductLineRequest struct {
	Code                            string  `json:"code" binding:"required"`
	Name                            string  `json:"name" binding:"required"`
	InventoryCOACode                *string `json:"inventory_coa_code"`
	CogsCOACode                     *string `json:"cogs_coa_code"`
	PurchaseDiscountCOACode         *string `json:"purchase_discount_coa_code"`
	PurchaseTaxCOACode              *string `json:"purchase_tax_coa_code"`
	AssetCOACode                    *string `json:"asset_coa_code"`
	AssetAccumulationCOACode        *string `json:"asset_accumulation_coa_code"`
	AssetDepreciationExpenseCOACode *string `json:"asset_depreciation_expense_coa_code"`
	IsActive                        *bool   `json:"is_active"`
}

type ProductLineResponse struct {
	ID                              string    `json:"id"`
	Code                            string    `json:"code"`
	Name                            string    `json:"name"`
	InventoryCOACode                *string   `json:"inventory_coa_code,omitempty"`
	CogsCOACode                     *string   `json:"cogs_coa_code,omitempty"`
	PurchaseDiscountCOACode         *string   `json:"purchase_discount_coa_code,omitempty"`
	PurchaseTaxCOACode              *string   `json:"purchase_tax_coa_code,omitempty"`
	AssetCOACode                    *string   `json:"asset_coa_code,omitempty"`
	AssetAccumulationCOACode        *string   `json:"asset_accumulation_coa_code,omitempty"`
	AssetDepreciationExpenseCOACode *string   `json:"asset_depreciation_expense_coa_code,omitempty"`
	IsActive                        bool      `json:"is_active"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

func ToProductLineResponse(pl *ItemProductLine) ProductLineResponse {
	return ProductLineResponse{
		ID:                              pl.ID,
		Code:                            pl.Code,
		Name:                            pl.Name,
		InventoryCOACode:                pl.InventoryCOACode,
		CogsCOACode:                     pl.CogsCOACode,
		PurchaseDiscountCOACode:         pl.PurchaseDiscountCOACode,
		PurchaseTaxCOACode:              pl.PurchaseTaxCOACode,
		AssetCOACode:                    pl.AssetCOACode,
		AssetAccumulationCOACode:        pl.AssetAccumulationCOACode,
		AssetDepreciationExpenseCOACode: pl.AssetDepreciationExpenseCOACode,
		IsActive:                        pl.IsActive,
		CreatedAt:                       pl.CreatedAt,
		UpdatedAt:                       pl.UpdatedAt,
	}
}

// ==========================================
// Product Line Service (Jalur A)
// ==========================================

type ProductLineService interface {
	Create(ctx context.Context, req CreateProductLineRequest, actor string) (*ProductLineResponse, error)
	Update(ctx context.Context, id string, req UpdateProductLineRequest, actor string) (*ProductLineResponse, error)
	Delete(ctx context.Context, id string, actor string) error
	GetByID(ctx context.Context, id string) (*ProductLineResponse, error)
	List(ctx context.Context, params ProductLineListParams) ([]ProductLineResponse, response.PaginationMeta, error)
}

type productLineService struct {
	repo Repository
}

func NewProductLineService(repo Repository) ProductLineService {
	return &productLineService{repo: repo}
}

func (s *productLineService) Create(ctx context.Context, req CreateProductLineRequest, actor string) (*ProductLineResponse, error) {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)

	if code == "" || name == "" {
		return nil, errors.New("kode dan nama lini produk wajib diisi")
	}

	existing, err := s.repo.FindProductLineByCode(ctx, code)
	if err == nil && existing != nil {
		return nil, ErrProductLineCodeAlreadyExists
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	now := time.Now()
	pl := ItemProductLine{
		ID:                              uuid.NewV7().String(),
		Code:                            code,
		Name:                            name,
		InventoryCOACode:                req.InventoryCOACode,
		CogsCOACode:                     req.CogsCOACode,
		PurchaseDiscountCOACode:         req.PurchaseDiscountCOACode,
		PurchaseTaxCOACode:              req.PurchaseTaxCOACode,
		AssetCOACode:                    req.AssetCOACode,
		AssetAccumulationCOACode:        req.AssetAccumulationCOACode,
		AssetDepreciationExpenseCOACode: req.AssetDepreciationExpenseCOACode,
		IsActive:                        isActive,
		CreatedAt:                       now,
		CreatedBy:                       actor,
		UpdatedAt:                       now,
		UpdatedBy:                       actor,
	}

	if err := s.repo.CreateProductLine(ctx, &pl); err != nil {
		return nil, err
	}

	res := ToProductLineResponse(&pl)
	return &res, nil
}

func (s *productLineService) Update(ctx context.Context, id string, req UpdateProductLineRequest, actor string) (*ProductLineResponse, error) {
	pl, err := s.repo.FindProductLineByID(ctx, id)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("kode dan nama lini produk wajib diisi")
	}

	if code != pl.Code {
		existing, err := s.repo.FindProductLineByCode(ctx, code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrProductLineCodeAlreadyExists
		}
	}

	if req.IsActive != nil {
		pl.IsActive = *req.IsActive
	}

	pl.Code = code
	pl.Name = name
	pl.InventoryCOACode = req.InventoryCOACode
	pl.CogsCOACode = req.CogsCOACode
	pl.PurchaseDiscountCOACode = req.PurchaseDiscountCOACode
	pl.PurchaseTaxCOACode = req.PurchaseTaxCOACode
	pl.AssetCOACode = req.AssetCOACode
	pl.AssetAccumulationCOACode = req.AssetAccumulationCOACode
	pl.AssetDepreciationExpenseCOACode = req.AssetDepreciationExpenseCOACode
	pl.UpdatedAt = time.Now()
	pl.UpdatedBy = actor

	if err := s.repo.UpdateProductLine(ctx, pl); err != nil {
		return nil, err
	}

	res := ToProductLineResponse(pl)
	return &res, nil
}

func (s *productLineService) Delete(ctx context.Context, id string, actor string) error {
	_, err := s.repo.FindProductLineByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteProductLine(ctx, id, actor)
}

func (s *productLineService) GetByID(ctx context.Context, id string) (*ProductLineResponse, error) {
	pl, err := s.repo.FindProductLineByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToProductLineResponse(pl)
	return &res, nil
}

func (s *productLineService) List(ctx context.Context, params ProductLineListParams) ([]ProductLineResponse, response.PaginationMeta, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	productLines, total, err := s.repo.FindAllProductLines(ctx, params)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	res := make([]ProductLineResponse, len(productLines))
	for i := range productLines {
		res[i] = ToProductLineResponse(&productLines[i])
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.Limit)))
	}

	meta := response.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return res, meta, nil
}

// ==========================================
// Product Line HTTP Handler
// ==========================================

type ProductLineHandler struct {
	service ProductLineService
}

func NewProductLineHandler(service ProductLineService) *ProductLineHandler {
	return &ProductLineHandler{service: service}
}

func (h *ProductLineHandler) RegisterRoutes(router *gin.RouterGroup) {
	productLines := router.Group("/item-product-lines")
	{
		productLines.GET("", middleware.RequirePermission("item:read"), h.List)
		productLines.GET("/:id", middleware.RequirePermission("item:read"), h.GetByID)
		productLines.POST("", middleware.RequirePermission("item:create"), h.Create)
		productLines.PUT("/:id", middleware.RequirePermission("item:update"), h.Update)
		productLines.DELETE("/:id", middleware.RequirePermission("item:delete"), h.Delete)
	}
}

func (h *ProductLineHandler) Create(c *gin.Context) {
	var req CreateProductLineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload lini produk gagal: "+err.Error(), nil)
		return
	}

	actor := getOperator(c)
	result, err := h.service.Create(c.Request.Context(), req, actor)
	if err != nil {
		if errors.Is(err, ErrProductLineCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Lini produk berhasil dibuat", result)
}

func (h *ProductLineHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProductLineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload lini produk gagal: "+err.Error(), nil)
		return
	}

	actor := getOperator(c)
	result, err := h.service.Update(c.Request.Context(), id, req, actor)
	if err != nil {
		if errors.Is(err, ErrProductLineNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrProductLineCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Lini produk berhasil diperbarui", result)
}

func (h *ProductLineHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	actor := getOperator(c)

	err := h.service.Delete(c.Request.Context(), id, actor)
	if err != nil {
		if errors.Is(err, ErrProductLineNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Lini produk berhasil dihapus", nil)
}

func (h *ProductLineHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	result, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductLineNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Detail lini produk berhasil dimuat", result)
}

type QueryProductLinesRequest struct {
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=10"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

func (h *ProductLineHandler) List(c *gin.Context) {
	var req QueryProductLinesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Parameter query tidak valid: "+err.Error(), nil)
		return
	}

	params := ProductLineListParams{
		Page:     req.Page,
		Limit:    req.Limit,
		Search:   req.Search,
		IsActive: req.IsActive,
	}

	productLines, meta, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Daftar lini produk berhasil dimuat", productLines, meta)
}
